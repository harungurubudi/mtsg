package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeContainer(t *testing.T) {
	// Test that we can initialize the container
	container := InitializeContainer()

	// Verify it's not nil
	assert.NotNil(t, container)

	// Simple type check - just verify it's not nil and can be called
	t.Logf("Container initialized successfully: %T", container)
}
