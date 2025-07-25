package user

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Email represents a user's email address.
type Email string

// Password represents a user's password.
type Password string

// Validate checks if the password meets the required specification:
// at least 8 characters, 1 uppercase, 1 lowercase, 1 number, 1 special character.
func (p Password) Validate() bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(p) < 8 {
		return false
	}
	for _, c := range p {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// CipherText represents a hashed password.
type CipherText string

// NewCipherText hashes the password using bcrypt (at least 10 rounds).
func NewCipherText(p Password) (CipherText, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	return CipherText(hash), err
}

// Matches checks if the provided password matches the ciphertext hash.
func (c CipherText) Matches(p Password) bool {
	return bcrypt.CompareHashAndPassword([]byte(c), []byte(p)) == nil
}
