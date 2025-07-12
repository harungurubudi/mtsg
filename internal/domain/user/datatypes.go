package user

import (
    "golang.org/x/crypto/bcrypt"
)

// Email represents a user's email address.
type Email string

// Password represents a user's password.
type Password string

// Ciphertext represents a hashed password.
type Ciphertext string

// NewCiphertext hashes the password using bcrypt (at least 10 rounds).
func NewCiphertext(p Password) (Ciphertext, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(p), 10)
    return Ciphertext(hash), err
}

// Matches checks if the provided password matches the ciphertext hash.
func (c Ciphertext) Matches(p Password) bool {
    return bcrypt.CompareHashAndPassword([]byte(c), []byte(p)) == nil
} 