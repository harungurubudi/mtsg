package error

// HTTPError represents controlled errors that map directly to an HTTP response.
// Should be used for validation errors, not found errors, unauthorized, etc.
// Does not contain stack trace.
type HTTPError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return e.Message
}

// GetStatusCode returns the HTTP status code
func (e *HTTPError) GetStatusCode() int {
	return e.StatusCode
}

// GetCode returns the machine-readable error code
func (e *HTTPError) GetCode() string {
	return e.Code
}

// NewHTTPError creates a new HTTPError with the given parameters
func NewHTTPError(statusCode int, code, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// NewNotFoundError creates a 404 Not Found error
func NewNotFoundError(message string) *HTTPError {
	return NewHTTPError(404, "not_found", message)
}

// NewValidationError creates a 400 Bad Request validation error
func NewValidationError(message string) *HTTPError {
	return NewHTTPError(400, "validation_error", message)
}

// NewUnauthorizedError creates a 401 Unauthorized error
func NewUnauthorizedError(message string) *HTTPError {
	return NewHTTPError(401, "unauthorized", message)
}

// NewForbiddenError creates a 403 Forbidden error
func NewForbiddenError(message string) *HTTPError {
	return NewHTTPError(403, "forbidden", message)
}

// NewInternalServerError creates a 500 Internal Server Error
func NewInternalServerError(message string) *HTTPError {
	return NewHTTPError(500, "internal_server_error", message)
}

// Common error instances for reuse
var (
	ErrNotFound     = NewNotFoundError("Data not found")
	ErrUnauthorized = NewUnauthorizedError("Unauthorized access")
	ErrForbidden    = NewForbiddenError("Access forbidden")
)
