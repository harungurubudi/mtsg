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

func TestErrorHandlerIntegration(t *testing.T) {
	// Setup Echo server with our error handler
	e := echo.New()

	config := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   false,
	}
	errorHandler := NewErrorHandler(config)
	e.HTTPErrorHandler = errorHandler.HandleError

	// Add a test route that returns an error
	e.GET("/test-error", func(c echo.Context) error {
		return pkgerror.NewUnauthorizedError("test unauthorized")
	})

	e.GET("/test-echo-error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "test bad request")
	})

	e.GET("/test-generic-error", func(c echo.Context) error {
		return assert.AnError
	})

	// Test HTTPError handling
	t.Run("HTTPError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-error", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "error", response["status"])
		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "unauthorized", errorData["code"])
		assert.Equal(t, "test unauthorized", errorData["message"])
		assert.Equal(t, float64(http.StatusUnauthorized), errorData["status_code"])
	})

	// Test Echo HTTPError handling
	t.Run("EchoHTTPError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-echo-error", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "error", response["status"])
		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "echo_error", errorData["code"])
		assert.Equal(t, "test bad request", errorData["message"])
		assert.Equal(t, float64(http.StatusBadRequest), errorData["status_code"])
	})

	// Test generic error handling
	t.Run("GenericError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-generic-error", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "error", response["status"])
		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "internal_server_error", errorData["code"])
		assert.Equal(t, "An internal error occurred", errorData["message"])
		assert.Equal(t, float64(http.StatusInternalServerError), errorData["status_code"])
		assert.Equal(t, assert.AnError.Error(), errorData["details"])
	})

	// Test 404 handling
	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "error", response["status"])
		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "echo_error", errorData["code"])
		assert.Equal(t, "Not Found", errorData["message"])
		assert.Equal(t, float64(http.StatusNotFound), errorData["status_code"])
	})
}
