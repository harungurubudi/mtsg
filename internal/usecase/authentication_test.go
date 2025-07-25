package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/harungurubudi/mtsg/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AuthenticationTestSuite provides a test suite for authentication usecase
type AuthenticationTestSuite struct {
	suite.Suite
}

// TestAuthentication runs the authentication test suite
func TestAuthentication(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

// TestLogin_Success tests successful login scenario
func (suite *AuthenticationTestSuite) TestLogin_Success() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test user with properly hashed password
	userID := user.UserID(uuid.New())
	tenantID := tenant.TenantID(uuid.New())
	password := user.Password("Password123!")
	ciphertext, _ := user.NewCipherText(password)

	testUser := &user.User{
		ID:         userID,
		TenantID:   tenantID,
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "admin",
		CipherText: ciphertext,
		Status:     user.UserStatusActive,
	}

	// Create test credential
	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: user.Password(password),
	}

	// Mock expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(testUser, nil)
	mockTokenGen.EXPECT().Generate(ctx, mock.AnythingOfType("token.Claims")).Return(token.Token("access_token"), nil).Times(3)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), session)
	assert.NotEmpty(suite.T(), session.AccessToken)
	assert.NotEmpty(suite.T(), session.RefreshToken)
}

// TestLogin_UserNotFound tests login when user is not found
func (suite *AuthenticationTestSuite) TestLogin_UserNotFound() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	credential := &authentication.Credential{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(nil, user.ErrUserNotFound)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), session)
	assert.Equal(suite.T(), authentication.ErrUserNotFound, err)
}

// TestLogin_InvalidPassword tests login with invalid password
func (suite *AuthenticationTestSuite) TestLogin_InvalidPassword() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test user with properly hashed password
	userID := user.UserID(uuid.New())
	tenantID := tenant.TenantID(uuid.New())
	password := user.Password("Password123!")
	ciphertext, _ := user.NewCipherText(password)

	testUser := &user.User{
		ID:         userID,
		TenantID:   tenantID,
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "admin",
		CipherText: ciphertext,
		Status:     user.UserStatusActive,
	}

	// Create test credential with wrong password
	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Mock expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(testUser, nil)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), session)
	assert.Equal(suite.T(), authentication.ErrInvalidCredential, err)
}

// TestLogin_UserInactive tests login with inactive user
func (suite *AuthenticationTestSuite) TestLogin_UserInactive() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test user with inactive status and properly hashed password
	userID := user.UserID(uuid.New())
	tenantID := tenant.TenantID(uuid.New())
	password := user.Password("Password123!")
	ciphertext, _ := user.NewCipherText(password)

	testUser := &user.User{
		ID:         userID,
		TenantID:   tenantID,
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "admin",
		CipherText: ciphertext,
		Status:     user.UserStatusInactive,
	}

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: password,
	}

	// Mock expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(testUser, nil)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), session)
	assert.Equal(suite.T(), authentication.ErrUserInactive, err)
}

// TestLogin_RepositoryError tests login when repository returns an error
func (suite *AuthenticationTestSuite) TestLogin_RepositoryError() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	repoError := errors.New("database connection failed")
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(nil, repoError)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), session)
	assert.Contains(suite.T(), err.Error(), "failed to find user by email")
}

// TestLogin_TokenGenerationError tests login when token generation fails
func (suite *AuthenticationTestSuite) TestLogin_TokenGenerationError() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test user with properly hashed password
	userID := user.UserID(uuid.New())
	tenantID := tenant.TenantID(uuid.New())
	password := user.Password("Password123!")
	ciphertext, _ := user.NewCipherText(password)

	testUser := &user.User{
		ID:         userID,
		TenantID:   tenantID,
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "admin",
		CipherText: ciphertext,
		Status:     user.UserStatusActive,
	}

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: password,
	}

	// Mock expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(testUser, nil)
	tokenError := errors.New("token generation failed")
	mockTokenGen.EXPECT().Generate(ctx, mock.AnythingOfType("token.Claims")).Return(token.Token(""), tokenError)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), session)
	assert.Contains(suite.T(), err.Error(), "failed to generate session")
}

// TestVerifyToken_Success tests successful token verification
func (suite *AuthenticationTestSuite) TestVerifyToken_Success() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test user
	userID := user.UserID(uuid.New())
	tenantID := tenant.TenantID(uuid.New())
	testUser := &user.User{
		ID:         userID,
		TenantID:   tenantID,
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "admin",
		CipherText: "hashedpassword",
		Status:     user.UserStatusActive,
	}

	// Create test token
	testToken := token.Token("valid_token")
	subject := authentication.AccessTokenSubject

	// Mock expectations
	expectedClaims := &token.Claims{
		Subject:    subject,
		Identifier: uuid.UUID(userID).String(),
		EXP:        1234567890,
		JTI:        uuid.New(),
	}
	mockTokenGen.EXPECT().Validate(ctx, testToken).Return(expectedClaims, nil)
	mockUserRepo.EXPECT().GetOneByID(ctx, userID).Return(testUser, nil)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	resultUser, err := authUsecase.VerifyToken(ctx, testToken, subject, tenantID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resultUser)
	assert.Equal(suite.T(), testUser.ID, resultUser.ID)
	assert.Equal(suite.T(), testUser.TenantID, resultUser.TenantID)
}

// TestVerifyToken_InvalidToken tests token verification with invalid token
func (suite *AuthenticationTestSuite) TestVerifyToken_InvalidToken() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test token
	testToken := token.Token("invalid_token")
	subject := authentication.AccessTokenSubject
	tenantID := tenant.TenantID(uuid.New())

	// Mock expectations
	tokenError := errors.New("token validation failed")
	mockTokenGen.EXPECT().Validate(ctx, testToken).Return((*token.Claims)(nil), tokenError)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	resultUser, err := authUsecase.VerifyToken(ctx, testToken, subject, tenantID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resultUser)
	assert.Equal(suite.T(), authentication.ErrInvalidAuthentication, err)
}

// TestVerifyToken_UserNotFound tests token verification when user is not found
func (suite *AuthenticationTestSuite) TestVerifyToken_UserNotFound() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test token
	testToken := token.Token("valid_token")
	subject := authentication.AccessTokenSubject
	tenantID := tenant.TenantID(uuid.New())
	userID := user.UserID(uuid.New())

	// Mock expectations
	expectedClaims := &token.Claims{
		Subject:    subject,
		Identifier: uuid.UUID(userID).String(),
		EXP:        1234567890,
		JTI:        uuid.New(),
	}
	mockTokenGen.EXPECT().Validate(ctx, testToken).Return(expectedClaims, nil)
	mockUserRepo.EXPECT().GetOneByID(ctx, userID).Return(nil, user.ErrUserNotFound)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	resultUser, err := authUsecase.VerifyToken(ctx, testToken, subject, tenantID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resultUser)
	assert.Equal(suite.T(), authentication.ErrInvalidAuthentication, err)
}

// TestVerifyToken_TenantMismatch tests token verification with tenant mismatch
func (suite *AuthenticationTestSuite) TestVerifyToken_TenantMismatch() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test user with different tenant
	userID := user.UserID(uuid.New())
	userTenantID := tenant.TenantID(uuid.New())
	requiredTenantID := tenant.TenantID(uuid.New()) // Different tenant
	testUser := &user.User{
		ID:         userID,
		TenantID:   userTenantID,
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "admin",
		CipherText: "hashedpassword",
		Status:     user.UserStatusActive,
	}

	// Create test token
	testToken := token.Token("valid_token")
	subject := authentication.AccessTokenSubject

	// Mock expectations
	expectedClaims := &token.Claims{
		Subject:    subject,
		Identifier: uuid.UUID(userID).String(),
		EXP:        1234567890,
		JTI:        uuid.New(),
	}
	mockTokenGen.EXPECT().Validate(ctx, testToken).Return(expectedClaims, nil)
	mockUserRepo.EXPECT().GetOneByID(ctx, userID).Return(testUser, nil)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	resultUser, err := authUsecase.VerifyToken(ctx, testToken, subject, requiredTenantID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resultUser)
	assert.Equal(suite.T(), authentication.ErrTenantMismatch, err)
}

// TestVerifyToken_RepositoryError tests token verification when repository returns an error
func (suite *AuthenticationTestSuite) TestVerifyToken_RepositoryError() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test token
	testToken := token.Token("valid_token")
	subject := authentication.AccessTokenSubject
	tenantID := tenant.TenantID(uuid.New())
	userID := user.UserID(uuid.New())

	// Mock expectations
	expectedClaims := &token.Claims{
		Subject:    subject,
		Identifier: uuid.UUID(userID).String(),
		EXP:        1234567890,
		JTI:        uuid.New(),
	}
	mockTokenGen.EXPECT().Validate(ctx, testToken).Return(expectedClaims, nil)
	repoError := errors.New("database connection failed")
	mockUserRepo.EXPECT().GetOneByID(ctx, userID).Return(nil, repoError)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	resultUser, err := authUsecase.VerifyToken(ctx, testToken, subject, tenantID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resultUser)
	assert.Contains(suite.T(), err.Error(), "failed to get user by ID")
}

// TestVerifyToken_InvalidUserID tests token verification with invalid user ID in claims
func (suite *AuthenticationTestSuite) TestVerifyToken_InvalidUserID() {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(suite.T())
	mockTokenGen := testmock.NewMockGeneratorRepository(suite.T())

	// Create test token
	testToken := token.Token("valid_token")
	subject := authentication.AccessTokenSubject
	tenantID := tenant.TenantID(uuid.New())

	// Mock expectations with invalid UUID
	expectedClaims := &token.Claims{
		Subject:    subject,
		Identifier: "invalid-uuid",
		EXP:        1234567890,
		JTI:        uuid.New(),
	}
	mockTokenGen.EXPECT().Validate(ctx, testToken).Return(expectedClaims, nil)

	// Create usecase
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	// Act
	resultUser, err := authUsecase.VerifyToken(ctx, testToken, subject, tenantID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resultUser)
	assert.Contains(suite.T(), err.Error(), "failed to parse user ID from token claims")
}
