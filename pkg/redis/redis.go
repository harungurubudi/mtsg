package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// Adapter defines the interface for the Redis adapter.
type AdapterRepository interface {
	// GetByKey retrieves the value for the given key and unmarshals/casts it into the provided receiver pointer.
	// Returns an error if the key does not exist or unmarshaling fails.
	GetByKey(ctx context.Context, key string, receiver any) error

	// Set stores the value under the given key with the specified expiration.
	// Value can be any serializable type.
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// DeleteByKeys deletes one or more keys from Redis. Returns an error if any deletion fails.
	DeleteByKeys(ctx context.Context, keys ...string) error

	// IsExist checks if the given key exists in Redis. Returns true if the key exists, false otherwise, and an error if the operation fails.
	IsExist(ctx context.Context, key string) (bool, error)
}

// adapter implements the Adapter interface and wraps a go-redis client.
type adapter struct {
	client *redis.Client
}

// NewAdapter creates a new Redis adapter with the given go-redis client.
func NewAdapter(client *redis.Client) AdapterRepository {
	return &adapter{client: client}
}

// GetByKey retrieves the value for the given key and unmarshals/casts it into the provided receiver pointer.
// Returns an error if the key does not exist or unmarshaling fails.
func (a *adapter) GetByKey(ctx context.Context, key string, receiver any) error {
	val, err := a.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), receiver)
}

// Set stores the value under the given key with the specified expiration.
// Value can be any serializable type.
func (a *adapter) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, data, expiration).Err()
}

// DeleteByKeys deletes one or more keys from Redis. Returns an error if any deletion fails.
func (a *adapter) DeleteByKeys(ctx context.Context, keys ...string) error {
	return a.client.Del(ctx, keys...).Err()
}

// IsExist checks if the given key exists in Redis. Returns true if the key exists, false otherwise, and an error if the operation fails.
func (a *adapter) IsExist(ctx context.Context, key string) (bool, error) {
	n, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
