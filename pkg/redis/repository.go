package redis

import (
    "context"
    "time"
)

// Adapter defines the interface for the Redis adapter.
type Adapter interface {
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