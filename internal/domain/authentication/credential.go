package authentication

import (
	"context"
	"errors"

	"github.com/harungurubudi/mtsg/internal/domain/user"
)

// Error types for specific authentication scenarios
var (
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUserNotFound      = errors.New("user not found")
	ErrAccountLocked     = errors.New("account is locked")
	ErrPasswordExpired   = errors.New("password has expired")
	ErrUserInactive      = errors.New("user account is inactive")
)

// Credential represents the input required for user authentication (login)
type Credential struct {
	Email    user.Email    `json:"email" validate:"required,email"`
	Password user.Password `json:"password" validate:"required,password"`
}

// Verify checks if the provided password matches the user's stored ciphertext.
// Returns nil if matched, otherwise returns ErrInvalidCredential.
func (c *Credential) Verify(ctx context.Context, matchedUser *user.User) error {
	if !c.Password.Validate() {
		return ErrInvalidCredential
	}
	// Convert CipherText to Ciphertext for the Matches method
	ciphertext := user.Ciphertext(matchedUser.CipherText)
	if !ciphertext.Matches(c.Password) {
		return ErrInvalidCredential
	}
	return nil
}
