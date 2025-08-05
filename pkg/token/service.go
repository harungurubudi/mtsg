package token

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/pkg/redis"
)

// Error definitions for token package.
var (
	ErrTokenNotFound    = errors.New("token not found or invalid")
	ErrTokenNotYetValid = errors.New("token not valid yet (nbf)")
)

type GeneratorRepository interface {
	Generate(ctx context.Context, claims Claims) (Token, error)
	Validate(ctx context.Context, token Token) (*Claims, error)
	Revoke(ctx context.Context, token Token) error
}

// Generator provides stateful token generation, validation, and revocation using Redis.
type generator struct {
	redis redis.AdapterRepository
	key   []byte
}

// NewGenerator creates a new token generator with the given Redis adapter and secret key.
func NewGenerator(redisAdapter redis.AdapterRepository, key string) GeneratorRepository {
	return &generator{
		redis: redisAdapter,
		key:   []byte(key),
	}
}

// Generate creates a new, unique, and unpredictable token, stores the claims in Redis, and returns the token.
func (g *generator) Generate(ctx context.Context, claims Claims) (Token, error) {
	// Generate a random nonce
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	// Use HMAC-SHA256 with the key and nonce to generate the token string
	mac := hmac.New(sha256.New, g.key)
	mac.Write(nonce)
	tokenBytes := mac.Sum(nil)
	tokenStr := base64.URLEncoding.EncodeToString(tokenBytes)
	token := Token(tokenStr)

	// Set JTI if not already set
	if claims.JTI == uuid.Nil {
		claims.JTI = uuid.New()
	}

	// Calculate TTL (expiration)
	ttl := time.Until(time.Unix(claims.EXP, 0))
	if ttl <= 0 {
		return "", errors.New("claims.EXP must be in the future")
	}

	// Store claims directly in Redis (Redis adapter handles JSON marshalling)
	if err := g.redis.Set(ctx, tokenStr, claims, ttl); err != nil {
		return "", fmt.Errorf("failed to store token in redis: %w", err)
	}

	return token, nil
}

// Validate checks if the token exists in Redis and returns the claims if valid.
func (g *generator) Validate(ctx context.Context, token Token) (*Claims, error) {
	var claims Claims
	err := g.redis.GetByKey(ctx, string(token), &claims)
	if err != nil {
		return nil, ErrTokenNotFound
	}
	now := time.Now().Unix()
	if claims.NBF > 0 && now < claims.NBF {
		return nil, ErrTokenNotYetValid
	}
	return &claims, nil
}

// Revoke deletes the token from Redis.
func (g *generator) Revoke(ctx context.Context, token Token) error {
	return g.redis.DeleteByKeys(ctx, string(token))
}
