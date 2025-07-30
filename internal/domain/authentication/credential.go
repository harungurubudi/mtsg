package authentication

import (
	"context"

	"github.com/harungurubudi/mtsg/internal/domain/user"
)

// Credential represents the input required for user authentication (login)
type Credential struct {
	Email    user.Email    `json:"email"`
	Password user.Password `json:"password"`
}

// Verify checks if the provided password matches the user's stored ciphertext.
// Returns nil if matched, otherwise returns ErrInvalidCredential.
func (c *Credential) Verify(ctx context.Context, matchedUser *user.User) error {
	if !c.Password.Validate() {
		return ErrInvalidCredential
	}
	// Convert CipherText to CipherText for the Matches method
	ciphertext := user.CipherText(matchedUser.CipherText)
	if !ciphertext.Matches(c.Password) {
		return ErrInvalidCredential
	}
	return nil
}
