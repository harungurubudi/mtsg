package authentication_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	userdomain "github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (suite *AuthenticationTestSuite) SetupTest() {
	suite.ctx = context.Background()
}

func (suite *AuthenticationTestSuite) TestCredential_Verify_ValidPassword() {
	// Create a hashed password
	password := "ValidPass1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Create a user with the hashed password
	user := &userdomain.User{
		ID:         userdomain.UserID(uuid.New()),
		TenantID:   tenant.TenantID(uuid.New()),
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "member",
		CipherText: userdomain.CipherText(hashedPassword),
		Status:     "active",
	}

	// Create credential with matching password
	cred := authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password(password),
	}

	// Verify should succeed
	err := cred.Verify(suite.ctx, user)
	suite.NoError(err)
}

func (suite *AuthenticationTestSuite) TestCredential_Verify_InvalidPassword() {
	// Create a hashed password
	password := "ValidPass1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Create a user with the hashed password
	user := &userdomain.User{
		ID:         userdomain.UserID(uuid.New()),
		TenantID:   tenant.TenantID(uuid.New()),
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "member",
		CipherText: userdomain.CipherText(hashedPassword),
		Status:     "active",
	}

	// Create credential with wrong password
	cred := authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password("WrongPassword1!"),
	}

	// Verify should fail
	err := cred.Verify(suite.ctx, user)
	suite.Error(err)
	suite.Equal(authentication.ErrInvalidCredential, err)
}

func (suite *AuthenticationTestSuite) TestCredential_Verify_InvalidPasswordFormat() {
	// Create a user
	user := &userdomain.User{
		ID:         userdomain.UserID(uuid.New()),
		TenantID:   tenant.TenantID(uuid.New()),
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       "member",
		CipherText: userdomain.CipherText("hashedpassword"),
		Status:     "active",
	}

	// Create credential with invalid password format
	cred := authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password(""), // Empty password should fail validation
	}

	// Verify should fail due to invalid password format
	err := cred.Verify(suite.ctx, user)
	suite.Error(err)
	suite.Equal(authentication.ErrInvalidCredential, err)
}

func (suite *AuthenticationTestSuite) TestSession_Creation() {
	session := authentication.Session{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
	}

	suite.Equal("access-token-123", session.AccessToken)
	suite.Equal("refresh-token-456", session.RefreshToken)
}

func (suite *AuthenticationTestSuite) TestErrorTypes() {
	// Test that error types are properly defined
	suite.Equal("invalid credentials", authentication.ErrInvalidCredential.Error())
	suite.Equal("user not found", authentication.ErrUserNotFound.Error())
	suite.Equal("account is locked", authentication.ErrAccountLocked.Error())
	suite.Equal("password has expired", authentication.ErrPasswordExpired.Error())
	suite.Equal("user account is inactive", authentication.ErrUserInactive.Error())
}

func (suite *AuthenticationTestSuite) TestTokenSubjectConstants() {
	// Test that token subject constants are properly defined
	suite.Equal("access_token", authentication.AccessTokenSubject)
	suite.Equal("refresh_token", authentication.RefreshTokenSubject)
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}
