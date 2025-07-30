package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware (will use nil auth use case for now)
	h := AuthMiddleware(nil)

	// Test
	err := h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	// Should return unauthorized error due to missing auth header
	assert.Error(t, err)
}

func TestCORSMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	h := CORSMiddleware()

	// Test
	err := h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	assert.NoError(t, err)
	// CORS middleware returns 204 for OPTIONS requests
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestLoggingMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	h := LoggingMiddleware()

	// Test
	err := h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRecoveryMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	h := RecoveryMiddleware()

	// Test with panic
	err := h(func(c echo.Context) error {
		panic("test panic")
	})(c)

	// Should recover from panic and return 500
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
	} else {
		// Recovery middleware should handle panic gracefully
		assert.NoError(t, err)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware (will use nil Redis adapter for now)
	h := RateLimitMiddleware(nil, nil)

	// Test
	err := h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	// Should pass through since Redis adapter is nil
	assert.NoError(t, err)
}

func TestValidationMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	h := ValidationMiddleware()

	// Test
	err := h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestContentTypeMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		allowedTypes   []string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "GET request should pass",
			method:         http.MethodGet,
			contentType:    "",
			allowedTypes:   []string{"application/json"},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "valid content type",
			method:         http.MethodPost,
			contentType:    "application/json",
			allowedTypes:   []string{"application/json"},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "invalid content type",
			method:         http.MethodPost,
			contentType:    "text/plain",
			allowedTypes:   []string{"application/json"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(tt.method, "/", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create middleware
			h := ContentTypeMiddleware(tt.allowedTypes...)

			// Test
			err := h(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})(c)

			if tt.expectedError {
				assert.Error(t, err)
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, httpErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestMiddlewareFactory(t *testing.T) {
	// Setup
	cfg := &config.Config{}

	// Create factory with nil dependencies for testing
	factory := NewMiddlewareFactory(nil, nil, cfg)

	// Test factory methods
	assert.NotNil(t, factory.GlobalMiddleware())
	assert.NotNil(t, factory.AuthMiddleware())
	assert.NotNil(t, factory.RateLimitMiddleware())
	assert.NotNil(t, factory.ValidationMiddleware())
	assert.NotNil(t, factory.ContentTypeMiddleware())
	assert.NotNil(t, factory.RecoveryMiddleware())
	assert.NotNil(t, factory.LoggingMiddleware())
	assert.NotNil(t, factory.CORSMiddleware())
}

func TestGetUserFromContext(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test with no user in context
	retrievedUser, exists := GetUserFromContext(c)
	assert.Nil(t, retrievedUser)
	assert.False(t, exists)

	// Test with user in context
	testUser := &user.User{}
	c.Set(UserContextKey, testUser)
	retrievedUser, exists = GetUserFromContext(c)
	assert.Equal(t, testUser, retrievedUser)
	assert.True(t, exists)
}

func TestRequireUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test without user in context
	h := RequireUser()
	err := h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	assert.Error(t, err)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	}

	// Test with user in context
	testUser := &user.User{}
	c.Set(UserContextKey, testUser)
	err = h(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
