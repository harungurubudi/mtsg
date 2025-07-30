package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeAuthUseCase(t *testing.T) {
	// Test that we can initialize the auth use case
	authUseCase := InitializeAuthUseCase()

	// Verify it's not nil
	assert.NotNil(t, authUseCase)

	// Simple type check - just verify it's not nil and can be called
	t.Logf("Auth use case initialized successfully: %T", authUseCase)
}
