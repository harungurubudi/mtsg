package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestServer_HealthCheck(t *testing.T) {
	// Create test server
	handlers := handler.NewHandlers()
	config := &config.Config{Server: config.ServerConfig{Port: "8080"}}
	server := NewServer(handlers, config)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := server.echo.NewContext(req, rec)

	// Test health check
	err := server.healthCheck(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify response contains expected fields
	responseBody := rec.Body.String()
	assert.Contains(t, responseBody, "status")
	assert.Contains(t, responseBody, "time")
}

func TestNewServer(t *testing.T) {
	// Test server creation
	handlers := handler.NewHandlers()
	config := &config.Config{Server: config.ServerConfig{Port: "8080"}}
	server := NewServer(handlers, config)

	// Verify server is not nil
	assert.NotNil(t, server)
	assert.NotNil(t, server.echo)
	assert.Equal(t, handlers, server.handlers)
	assert.Equal(t, config, server.config)
}

func TestNewConfig(t *testing.T) {
	// Test config creation
	config := NewConfig()

	// Verify config has sensible defaults
	assert.NotEmpty(t, config.Server.Port)
	assert.NotZero(t, config.Server.ReadTimeout)
	assert.NotZero(t, config.Server.WriteTimeout)
	assert.NotZero(t, config.Server.IdleTimeout)
}
