package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthentication is a simple mock for testing
type MockAuthentication struct {
	mock.Mock
}

func (m *MockAuthentication) Login(ctx context.Context, credential *authentication.Credential) (*authentication.Session, error) {
	return nil, nil
}

func (m *MockAuthentication) VerifyToken(ctx context.Context, token token.Token, subject string, tenantID tenant.TenantID) (*user.User, error) {
	return nil, nil
}

func (m *MockAuthentication) RefreshToken(ctx context.Context, refreshToken token.Token) (*authentication.Session, error) {
	return nil, nil
}

func (m *MockAuthentication) Logout(ctx context.Context, accessToken token.Token) error {
	return nil
}

func TestServer_HealthCheck(t *testing.T) {
	// Create test server
	mockAuth := new(MockAuthentication)
	handlers := handler.NewHandlers(mockAuth)
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
	mockAuth := new(MockAuthentication)
	handlers := handler.NewHandlers(mockAuth)
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
