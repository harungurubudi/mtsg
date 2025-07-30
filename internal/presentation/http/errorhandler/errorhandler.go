package errorhandler

import (
	"net/http"

	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/labstack/echo/v4"
)

// ErrorHandler provides centralized error handling
type ErrorHandler struct {
	config *Config
}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler(config *Config) *ErrorHandler {
	if config == nil {
		config = DefaultConfig()
	}

	return &ErrorHandler{
		config: config,
	}
}

// HandleError processes all types of errors and returns appropriate responses
func (h *ErrorHandler) HandleError(err error, c echo.Context) {
	// Log error if enabled
	if h.config.LogErrors {
		h.logError(err, c)
	}

	// Handle different error types
	switch e := err.(type) {
	case *pkgerror.HTTPError:
		h.handleHTTPError(e, c)
	case *pkgerror.StackError:
		h.handleStackError(e, c)
	case *echo.HTTPError:
		h.handleEchoHTTPError(e, c)
	default:
		h.handleGenericError(err, c)
	}
}

// handleHTTPError processes HTTPError from pkg/error
func (h *ErrorHandler) handleHTTPError(err *pkgerror.HTTPError, c echo.Context) {
	response := map[string]interface{}{
		"status": "error",
		"error": map[string]interface{}{
			"code":        err.GetCode(),
			"message":     err.Error(),
			"status_code": err.GetStatusCode(),
		},
	}

	c.JSON(err.GetStatusCode(), response)
}

// handleStackError processes StackError from pkg/error
func (h *ErrorHandler) handleStackError(err *pkgerror.StackError, c echo.Context) {
	response := map[string]interface{}{
		"status": "error",
		"error": map[string]interface{}{
			"code":        "internal_server_error",
			"message":     "An internal error occurred",
			"status_code": http.StatusInternalServerError,
		},
	}

	// Add stack trace in development
	if h.config.IsDevelopment() && h.config.ShowDetails {
		response["error"].(map[string]interface{})["stack_trace"] = err.StackTrace()
	}

	c.JSON(http.StatusInternalServerError, response)
}

// handleEchoHTTPError processes Echo's HTTPError
func (h *ErrorHandler) handleEchoHTTPError(err *echo.HTTPError, c echo.Context) {
	response := map[string]interface{}{
		"status": "error",
		"error": map[string]interface{}{
			"code":        "echo_error",
			"message":     err.Message,
			"status_code": err.Code,
		},
	}

	// Add internal details in development
	if h.config.IsDevelopment() && h.config.ShowDetails {
		response["error"].(map[string]interface{})["internal"] = err.Internal
	}

	c.JSON(err.Code, response)
}

// handleGenericError processes generic errors
func (h *ErrorHandler) handleGenericError(err error, c echo.Context) {
	response := map[string]interface{}{
		"status": "error",
		"error": map[string]interface{}{
			"code":        "internal_server_error",
			"message":     "An internal error occurred",
			"status_code": http.StatusInternalServerError,
		},
	}

	// Add error details in development
	if h.config.IsDevelopment() && h.config.ShowDetails {
		response["error"].(map[string]interface{})["details"] = err.Error()
	}

	c.JSON(http.StatusInternalServerError, response)
}

// logError logs error details for debugging
func (h *ErrorHandler) logError(err error, c echo.Context) {
	logger := c.Logger()

	logData := map[string]interface{}{
		"error":      err.Error(),
		"method":     c.Request().Method,
		"uri":        c.Request().RequestURI,
		"remote_ip":  c.RealIP(),
		"user_agent": c.Request().UserAgent(),
	}

	// Add stack trace for stack errors
	if stackErr, ok := err.(*pkgerror.StackError); ok {
		logData["stack_trace"] = stackErr.StackTrace()
	}

	logger.Error(logData)
}
