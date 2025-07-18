package token

import (
    "github.com/google/uuid"
)

// Claims represents the data inside a token.
type Claims struct {
    Subject     string                 `json:"sub"`                 // What the token is for
    Identifier  string                 `json:"id"`                  // Who the token belongs to
    Payload     map[string]interface{} `json:"payload,omitempty"`    // Meta information or custom claims
    LinkedToken Token                  `json:"linked_token,omitempty"` // Token to be revoked together
    JTI         uuid.UUID              `json:"jti"`                 // Unique identifier for the token (JWT ID)
    NBF         int64                  `json:"nbf"`                 // Not before (Unix seconds)
    EXP         int64                  `json:"exp"`                 // Expiration (Unix seconds)
}

// Token is a string representation of the claims (e.g., JWT, opaque string).
type Token string

// New creates a new token from claims.
// JTI should be generated with uuid.New() when creating a new token.
// Implementation should serialize, sign, and return the token string.
func New(claims Claims) (Token, error) {
    // TODO: Implement token serialization and signing
    return "", nil
}

// Validate parses and validates a token string.
// Implementation should parse, verify signature, check expiration, etc.
func Validate(token Token) (*Claims, error) {
    // TODO: Implement token parsing and validation
    return nil, nil
} 