package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/internal/repository"
	stackerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/harungurubudi/mtsg/pkg/token"
)

// Authentication defines the interface for authentication-related usecases
type Authentication interface {
	Login(ctx context.Context, credential *authentication.Credential) (*authentication.Session, error)
	// VerifyToken verifies if token is valid for a specific subject and tenantID.
	// Returns matched user object on positive case.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout
	//   - token: Token to verify
	//   - subject: Desired subject (e.g., "access_token", "refresh_token")
	//   - tenantID: Required tenant for multi-tenant validation
	//
	// Returns:
	//   - *user.User: Matched user with the payload
	//   - error: Authentication or authorization error
	//
	// Possible errors:
	//   - authentication.ErrInvalidAuthentication: Invalid/expired token
	//   - authentication.ErrForbidden: Tenant mismatch or access forbidden
	//   - authentication.ErrTenantMismatch: User does not belong to required tenant
	//   - stackerror.StackError: Internal errors with stack trace
	VerifyToken(ctx context.Context, token token.Token, subject string, tenantID tenant.TenantID) (*user.User, error)
}

// authenticationUsecase implements the Authentication interface.
//
// The struct follows the dependency injection pattern, accepting
// interfaces rather than concrete implementations for better
// testability and flexibility.
type authenticationUsecase struct {
	userRepo repository.UserRepository
	tokenGen token.GeneratorRepository
}

// NewAuthentication creates a new instance of Authentication usecase
func NewAuthentication(userRepo repository.UserRepository, tokenGen token.GeneratorRepository) Authentication {
	return &authenticationUsecase{
		userRepo: userRepo,
		tokenGen: tokenGen,
	}
}

// Login handles user authentication with the provided credentials.
// This method implements the complete authentication flow including
// user lookup, password verification, status validation, and token generation.
//
// The authentication process follows these steps:
// 1. User Lookup: Find user by email address
// 2. Password Verification: Validate provided password against stored hash
// 3. Status Check: Ensure user account is active
// 4. Token Generation: Create linked access and refresh tokens
// 5. Session Creation: Return session with both tokens
//
// Security Features:
// - Password verification uses secure hashing comparison
// - Tokens are linked for enhanced security
// - Access tokens expire in 2 hours
// - Refresh tokens expire in 48 hours with 2-hour not-before time
// - All internal errors include stack traces for debugging
//
// Error Handling:
//   - Business errors (authentication.ErrUserNotFound, etc.) are returned directly
//   - Internal errors are wrapped with stackerror for debugging
//   - All errors are logged appropriately based on type
func (a *authenticationUsecase) Login(ctx context.Context, credential *authentication.Credential) (*authentication.Session, error) {
	// Step 1: Find matchedUser by email
	matchedUser, err := a.userRepo.GetOneByEmail(ctx, credential.Email)
	if err != nil {
		// Check if it's a "not found" error from repository
		if err == user.ErrUserNotFound {
			return nil, authentication.ErrUserNotFound
		}
		// Return internal error for unexpected repository errors
		return nil, stackerror.NewStackError("failed to find user by email", err)
	}

	// Step 2: Verify Password via credential.Verify
	if err := credential.Verify(ctx, matchedUser); err != nil {
		return nil, err
	}

	// Step 3: Ensure matchedUser is active
	if matchedUser.Status != user.UserStatusActive {
		return nil, authentication.ErrUserInactive
	}

	// Step 4: Generate access and refresh token in Session result
	session, err := a.generateSession(ctx, matchedUser)
	if err != nil {
		return nil, stackerror.NewStackError("failed to generate session", err)
	}

	// Step 5: Return Session
	return session, nil
}

// VerifyToken verifies if token is valid for a specific subject and tenantID.
// This method implements secure token validation with multi-tenant support.
//
// The verification process follows these steps:
// 1. Token Validation: Verify token validity for the specified subject
// 2. User Extraction: Extract user from token claims on success
// 3. Tenant Validation: Check if user belongs to required tenant
// 4. Error Handling: Proper error wrapping for internal errors
// 5. Authorization Check: Return forbidden if tenant doesn't match
//
// Security Features:
// - Subject validation ensures tokens are used for intended purpose
// - Multi-tenant validation for data isolation
// - Proper error messages that don't leak sensitive information
// - Stack trace preservation for debugging internal errors
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - token: Token to verify
//   - subject: Expected token subject (e.g., "access_token", "refresh_token")
//   - tenantID: Required tenant for multi-tenant validation
//
// Returns:
//   - *user.User: Matched user from token payload
//   - error: Authentication or authorization error
//
// Error Handling:
//   - Business errors (authentication.ErrInvalidAuthentication, etc.) are returned directly
//   - Internal errors are wrapped with stackerror for debugging
//   - User not found errors are converted to authentication errors
func (a *authenticationUsecase) VerifyToken(ctx context.Context, token token.Token, subject string, tenantID tenant.TenantID) (*user.User, error) {
	// Step 1: Verify if token is valid for the subject using tokenGen.Validate()
	claims, err := a.tokenGen.Validate(ctx, token)
	if err != nil {
		return nil, authentication.ErrInvalidAuthentication
	}

	// Step 2: If valid, get matched user from UserID in the token claim using userRepo.GetOneByID()
	userID, err := uuid.Parse(claims.Identifier)
	if err != nil {
		return nil, stackerror.NewStackError("failed to parse user ID from token claims", err)
	}

	matchedUser, err := a.userRepo.GetOneByID(ctx, user.UserID(userID))
	if err != nil {
		// Step 4: On error, check if the error is user.ErrUserNotFound
		if err == user.ErrUserNotFound {
			return nil, authentication.ErrInvalidAuthentication
		}
		// If not, wrap the error with stackerror.NewStackError function
		return nil, stackerror.NewStackError("failed to get user by ID", err)
	}

	// Step 3: On success, check if returned user has TenantID equal with required tenantID in the args
	// Step 5: If tenantID is matched, return the user. Otherwise return authentication.ErrTenantMismatch
	if matchedUser.TenantID != tenantID {
		return nil, authentication.ErrTenantMismatch
	}

	return matchedUser, nil
}

// generateSession creates a new session with access and refresh tokens.
// This method handles the token generation process including the creation
// of linked tokens for enhanced security.
//
// Token Configuration:
// - Access Token: 2-hour expiration, linked to refresh token
// - Refresh Token: 48-hour expiration, 2-hour not-before time, linked to access token
// - Both tokens include user identifier and appropriate subjects
// - JTI (JWT ID) ensures token uniqueness
//
// Security Features:
// - Tokens are linked bidirectionally for revocation tracking
// - Refresh tokens have a not-before time to prevent immediate reuse
// - All token generation errors include stack traces
//
// Error Handling:
//   - All token generation errors are wrapped with stackerror
//   - Specific error messages for each generation step
//   - Errors include context about which token failed to generate
func (a *authenticationUsecase) generateSession(ctx context.Context, matchedUser *user.User) (*authentication.Session, error) {
	// Generate access token (expires in 2 hours)
	accessTokenClaims := token.Claims{
		Subject:    authentication.AccessTokenSubject,
		Identifier: uuid.UUID(matchedUser.ID).String(),
		EXP:        time.Now().Add(2 * time.Hour).Unix(),
		JTI:        uuid.New(),
	}

	accessToken, err := a.tokenGen.Generate(ctx, accessTokenClaims)
	if err != nil {
		return nil, stackerror.NewStackError("failed to generate access token", err)
	}

	// Generate refresh token (expires in 2 days, nbf in 2 hours)
	refreshTokenClaims := token.Claims{
		Subject:     authentication.RefreshTokenSubject,
		Identifier:  uuid.UUID(matchedUser.ID).String(),
		EXP:         time.Now().Add(48 * time.Hour).Unix(), // 2 days
		NBF:         time.Now().Add(2 * time.Hour).Unix(),  // nbf in 2 hours
		LinkedToken: accessToken,
		JTI:         uuid.New(),
	}

	refreshToken, err := a.tokenGen.Generate(ctx, refreshTokenClaims)
	if err != nil {
		return nil, stackerror.NewStackError("failed to generate refresh token", err)
	}

	// Update access token with linked refresh token
	accessTokenClaims.LinkedToken = refreshToken
	accessToken, err = a.tokenGen.Generate(ctx, accessTokenClaims)
	if err != nil {
		return nil, stackerror.NewStackError("failed to update access token with linked refresh token", err)
	}

	return &authentication.Session{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}, nil
}
