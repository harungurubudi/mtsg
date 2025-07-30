package errorhandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorHandler(t *testing.T) {
	// Test with nil config
	handler := NewErrorHandler(nil)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.config)
	assert.Equal(t, "development", handler.config.Environment)
	assert.True(t, handler.config.ShowDetails)
	assert.True(t, handler.config.LogErrors)

	// Test with custom config
	customConfig := &Config{
		Environment: "production",
		ShowDetails: false,
		LogErrors:   false,
	}
	handler = NewErrorHandler(customConfig)
	assert.NotNil(t, handler)
	assert.Equal(t, customConfig, handler.config)
}

func TestHandleHTTPError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   false,
	}
	handler := NewErrorHandler(config)

	// Test HTTPError handling
	httpErr := pkgerror.NewUnauthorizedError("unauthorized access")
	handler.HandleError(httpErr, c)

	// Verify response
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error", response["status"])
	errorData := response["error"].(map[string]interface{})
	assert.Equal(t, "unauthorized", errorData["code"])
	assert.Equal(t, "unauthorized access", errorData["message"])
	assert.Equal(t, float64(http.StatusUnauthorized), errorData["status_code"])
}

func TestHandleEchoHTTPError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   false,
	}
	handler := NewErrorHandler(config)

	// Test Echo HTTPError handling
	echoErr := echo.NewHTTPError(http.StatusBadRequest, "bad request")
	handler.HandleError(echoErr, c)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error", response["status"])
	errorData := response["error"].(map[string]interface{})
	assert.Equal(t, "echo_error", errorData["code"])
	assert.Equal(t, "bad request", errorData["message"])
	assert.Equal(t, float64(http.StatusBadRequest), errorData["status_code"])
}

func TestHandleGenericError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   false,
	}
	handler := NewErrorHandler(config)

	// Test generic error handling
	genericErr := assert.AnError
	handler.HandleError(genericErr, c)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error", response["status"])
	errorData := response["error"].(map[string]interface{})
	assert.Equal(t, "internal_server_error", errorData["code"])
	assert.Equal(t, "An internal error occurred", errorData["message"])
	assert.Equal(t, float64(http.StatusInternalServerError), errorData["status_code"])
	assert.Equal(t, genericErr.Error(), errorData["details"])
}

func TestHandleGenericErrorProduction(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := &Config{
		Environment: "production",
		ShowDetails: false,
		LogErrors:   false,
	}
	handler := NewErrorHandler(config)

	// Test generic error handling in production
	genericErr := assert.AnError
	handler.HandleError(genericErr, c)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error", response["status"])
	errorData := response["error"].(map[string]interface{})
	assert.Equal(t, "internal_server_error", errorData["code"])
	assert.Equal(t, "An internal error occurred", errorData["message"])
	assert.Equal(t, float64(http.StatusInternalServerError), errorData["status_code"])

	// Should not have details in production
	_, hasDetails := errorData["details"]
	assert.False(t, hasDetails)
}

func TestConfigMethods(t *testing.T) {
	// Test development config
	devConfig := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   true,
	}
	assert.True(t, devConfig.IsDevelopment())
	assert.False(t, devConfig.IsProduction())

	// Test production config
	prodConfig := &Config{
		Environment: "production",
		ShowDetails: false,
		LogErrors:   true,
	}
	assert.False(t, prodConfig.IsDevelopment())
	assert.True(t, prodConfig.IsProduction())
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "development", config.Environment)
	assert.True(t, config.ShowDetails)
	assert.True(t, config.LogErrors)
}

func TestErrorHandlingMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   false,
	}
	errorHandler := NewErrorHandler(config)
	middleware := ErrorHandlingMiddleware(errorHandler)

	// Test middleware with error
	err := middleware(func(c echo.Context) error {
		return pkgerror.NewNotFoundError("not found")
	})(c)

	// Should return nil since error is handled
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDevelopmentErrorHandlingMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := DevelopmentErrorHandlingMiddleware()

	// Test middleware with error
	err := middleware(func(c echo.Context) error {
		return assert.AnError
	})(c)

	// Should return nil since error is handled
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	jsonErr := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, jsonErr)

	errorData := response["error"].(map[string]interface{})
	_, hasDetails := errorData["details"]
	assert.True(t, hasDetails) // Should have details in development
}

func TestProductionErrorHandlingMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := ProductionErrorHandlingMiddleware()

	// Test middleware with error
	err := middleware(func(c echo.Context) error {
		return assert.AnError
	})(c)

	// Should return nil since error is handled
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	jsonErr := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, jsonErr)

	errorData := response["error"].(map[string]interface{})
	_, hasDetails := errorData["details"]
	assert.False(t, hasDetails) // Should not have details in production
}
