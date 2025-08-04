package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/user"
)

// LoginRequest represents login request data
// @Description Login credentials for user authentication
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com" description:"User's email address"`
	Password string `json:"password" validate:"required,min=8" example:"SecurePass123!" description:"User's password (minimum 8 characters)"`
}

// ToDomain converts LoginRequest to domain Credential
func (r *LoginRequest) ToDomain() authentication.Credential {
	return authentication.Credential{
		Email:    user.Email(r.Email),
		Password: user.Password(r.Password),
	}
}

// LoginResponse represents login response data
// @Description Authentication tokens returned after successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT access token for API authentication"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT refresh token for token renewal"`
}

// NewLoginResponse creates a new LoginResponse from domain Session
func NewLoginResponse(session *authentication.Session) *LoginResponse {
	if session == nil {
		return &LoginResponse{
			AccessToken:  "",
			RefreshToken: "",
		}
	}

	return &LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}
}

// RefreshTokenRequest represents refresh token request data
// @Description Request to refresh access token using refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"Valid refresh token for token renewal"`
}

// ToDomain converts RefreshTokenRequest to domain string
func (r *RefreshTokenRequest) ToDomain() string {
	return r.RefreshToken
}

// RegisterRequest represents user registration request data
// @Description User registration data (for future implementation)
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com" description:"User's email address"`
	Password string `json:"password" validate:"required,min=8" example:"SecurePass123!" description:"User's password (minimum 8 characters)"`
	Name     string `json:"name" validate:"required,min=1,max=100" example:"John Doe" description:"User's full name"`
}

// ToDomain converts RegisterRequest to domain User
func (r *RegisterRequest) ToDomain() user.User {
	return user.User{
		Email:      user.Email(r.Email),
		CipherText: user.CipherText(r.Password), // This will be hashed in the use case
		Name:       r.Name,
	}
}

// UserResponse represents user response data
// @Description User information returned in responses
type UserResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Unique user identifier"`
	Email     string `json:"email" example:"user@example.com" description:"User's email address"`
	Name      string `json:"name" example:"John Doe" description:"User's full name"`
	Role      string `json:"role" example:"member" description:"User's role (admin or member)"`
	Status    string `json:"status" example:"active" description:"User's status (active or inactive)"`
	CreatedAt string `json:"created_at" example:"2023-07-29T17:30:00Z" description:"User creation timestamp"`
	UpdatedAt string `json:"updated_at" example:"2023-07-29T17:30:00Z" description:"User last update timestamp"`
}

// NewUserResponse creates a new UserResponse from domain User
func NewUserResponse(user *user.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        uuid.UUID(user.ID).String(),
		Email:     string(user.Email),
		Name:      user.Name,
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
