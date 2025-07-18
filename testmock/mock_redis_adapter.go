package testmock

import (
    "context"
    "time"

    "github.com/stretchr/testify/mock"
)

// MockRedisAdapter is a testify/mock implementation of the redis.Adapter interface for testing.
type MockRedisAdapter struct {
    mock.Mock
}

// GetByKey mocks retrieving a value by key and unmarshaling it into the receiver.
func (m *MockRedisAdapter) GetByKey(ctx context.Context, key string, receiver any) error {
    args := m.Called(ctx, key, receiver)
    return args.Error(0)
}

// Set mocks storing a value under a key with expiration.
func (m *MockRedisAdapter) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
    args := m.Called(ctx, key, value, expiration)
    return args.Error(0)
}

// DeleteByKeys mocks deleting one or more keys from Redis.
func (m *MockRedisAdapter) DeleteByKeys(ctx context.Context, keys ...string) error {
    args := m.Called(ctx, keys)
    return args.Error(0)
}

// IsExist mocks checking if a key exists in Redis.
func (m *MockRedisAdapter) IsExist(ctx context.Context, key string) (bool, error) {
    args := m.Called(ctx, key)
    return args.Bool(0), args.Error(1)
} 