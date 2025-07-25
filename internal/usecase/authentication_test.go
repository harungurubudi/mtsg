package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	userdomain "github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/harungurubudi/mtsg/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(t)
	mockTokenGen := testmock.NewMockGeneratorRepository(t)

	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	password := "ValidPass1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &userdomain.User{
		ID:         userdomain.UserID(uuid.New()),
		TenantID:   tenant.TenantID(uuid.New()),
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       userdomain.UserRoleMember,
		CipherText: userdomain.CipherText(hashedPassword),
		Status:     userdomain.UserStatusActive,
	}

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password(password),
	}

	accessToken := token.Token("access-token-123")
	refreshToken := token.Token("refresh-token-456")

	// Expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(user, nil)
	mockTokenGen.EXPECT().Generate(mock.Anything, mock.Anything).Return(accessToken, nil).Once()
	mockTokenGen.EXPECT().Generate(mock.Anything, mock.Anything).Return(refreshToken, nil).Once()
	mockTokenGen.EXPECT().Generate(mock.Anything, mock.Anything).Return(accessToken, nil).Once()

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "access-token-123", session.AccessToken)
	assert.Equal(t, "refresh-token-456", session.RefreshToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(t)
	mockTokenGen := testmock.NewMockGeneratorRepository(t)
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	credential := &authentication.Credential{
		Email:    "nonexistent@example.com",
		Password: userdomain.Password("password"),
	}

	// Expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(nil, userdomain.ErrUserNotFound)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, authentication.ErrUserNotFound, err)
}

func TestLogin_InvalidPassword(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(t)
	mockTokenGen := testmock.NewMockGeneratorRepository(t)
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	password := "ValidPass1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &userdomain.User{
		ID:         userdomain.UserID(uuid.New()),
		TenantID:   tenant.TenantID(uuid.New()),
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       userdomain.UserRoleMember,
		CipherText: userdomain.CipherText(hashedPassword),
		Status:     userdomain.UserStatusActive,
	}

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password("WrongPassword"),
	}

	// Expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(user, nil)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, authentication.ErrInvalidCredential, err)
}

func TestLogin_UserInactive(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(t)
	mockTokenGen := testmock.NewMockGeneratorRepository(t)
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	password := "ValidPass1!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &userdomain.User{
		ID:         userdomain.UserID(uuid.New()),
		TenantID:   tenant.TenantID(uuid.New()),
		Email:      "test@example.com",
		Name:       "Test User",
		Role:       userdomain.UserRoleMember,
		CipherText: userdomain.CipherText(hashedPassword),
		Status:     userdomain.UserStatusInactive,
	}

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password(password),
	}

	// Expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(user, nil)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Equal(t, authentication.ErrUserInactive, err)
}

func TestLogin_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := testmock.NewMockUserRepository(t)
	mockTokenGen := testmock.NewMockGeneratorRepository(t)
	authUsecase := NewAuthentication(mockUserRepo, mockTokenGen)

	credential := &authentication.Credential{
		Email:    "test@example.com",
		Password: userdomain.Password("password"),
	}

	repoError := errors.New("database connection failed")

	// Expectations
	mockUserRepo.EXPECT().GetOneByEmail(ctx, credential.Email).Return(nil, repoError)

	// Act
	session, err := authUsecase.Login(ctx, credential)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "failed to find user by email")
}
