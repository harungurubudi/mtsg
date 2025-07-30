package dto

import (
	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/user"
)

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// ToDomain converts LoginRequest to domain Credential
func (r *LoginRequest) ToDomain() authentication.Credential {
	return authentication.Credential{
		Email:    user.Email(r.Email),
		Password: user.Password(r.Password),
	}
}

// LoginResponse represents login response data
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
