# Error Handler Package

This package provides centralized error handling for the Echo HTTP server, supporting multiple error types and environment-specific behavior.

## Features

- **Centralized Error Handling**: Single point of error processing for all HTTP errors
- **Multiple Error Types**: Support for HTTPError, StackError, Echo HTTPError, and generic errors
- **Environment-Specific Behavior**: Different error responses for development and production
- **Structured Logging**: Comprehensive error logging with request context
- **Factory Pattern**: Easy creation and configuration of error handlers
- **Middleware Support**: Error handling middleware for Echo middleware chains

## Directory Structure

```
internal/presentation/http/errorhandler/
├── config.go                # Error handler configuration
├── errorhandler.go          # Main error handler implementation
├── factory.go               # Error handler factory
├── middleware.go            # Error handling middleware
├── errorhandler_test.go     # Unit tests
├── integration_test.go      # Integration tests
└── README.md               # This file
```

## Usage

### Basic Usage

```go
import (
    "github.com/harungurubudi/mtsg/internal/presentation/http/errorhandler"
)

// Create error handler with default configuration
handler := errorhandler.NewErrorHandler(nil)

// Use with Echo server
e := echo.New()
e.HTTPErrorHandler = handler.HandleError
```

### Configuration-Based Usage

```go
// Create custom configuration
config := &errorhandler.Config{
    Environment: "development",
    ShowDetails: true,
    LogErrors:   true,
}

// Create error handler with configuration
handler := errorhandler.NewErrorHandler(config)
```

### Factory Pattern Usage

```go
import (
    "github.com/harungurubudi/mtsg/pkg/config"
)

// Create factory with application config
factory := errorhandler.NewErrorHandlerFactory(appConfig)

// Create error handler
handler := factory.CreateErrorHandler()

// Or create environment-specific handlers
devHandler := factory.CreateDevelopmentErrorHandler()
prodHandler := factory.CreateProductionErrorHandler()
```

### Middleware Usage

```go
// Create error handling middleware
middleware := errorhandler.ErrorHandlingMiddleware(handler)

// Use in Echo middleware chain
e.Use(middleware)

// Or use environment-specific middleware
e.Use(errorhandler.DevelopmentErrorHandlingMiddleware())
e.Use(errorhandler.ProductionErrorHandlingMiddleware())
```

## Error Types

### 1. HTTPError (pkg/error.HTTPError)

Custom HTTP errors from the `pkg/error` package.

```go
// Example response
{
    "status": "error",
    "error": {
        "code": "unauthorized",
        "message": "unauthorized access",
        "status_code": 401
    }
}
```

### 2. StackError (pkg/error.StackError)

Stack trace errors with debugging information.

```go
// Development response
{
    "status": "error",
    "error": {
        "code": "internal_server_error",
        "message": "An internal error occurred",
        "status_code": 500,
        "stack_trace": "stack trace details..."
    }
}

// Production response
{
    "status": "error",
    "error": {
        "code": "internal_server_error",
        "message": "An internal error occurred",
        "status_code": 500
    }
}
```

### 3. Echo HTTPError

Echo framework HTTP errors.

```go
// Example response
{
    "status": "error",
    "error": {
        "code": "echo_error",
        "message": "bad request",
        "status_code": 400
    }
}
```

### 4. Generic Error

Fallback error handling for unknown error types.

```go
// Development response
{
    "status": "error",
    "error": {
        "code": "internal_server_error",
        "message": "An internal error occurred",
        "status_code": 500,
        "details": "actual error message"
    }
}

// Production response
{
    "status": "error",
    "error": {
        "code": "internal_server_error",
        "message": "An internal error occurred",
        "status_code": 500
    }
}
```

## Configuration

### Config Structure

```go
type Config struct {
    Environment string `json:"environment"` // "development" or "production"
    ShowDetails bool   `json:"show_details"` // Show error details
    LogErrors   bool   `json:"log_errors"`   // Enable error logging
}
```

### Environment-Specific Behavior

#### Development Environment
- Shows detailed error messages
- Includes stack traces for StackError
- Logs all errors with full context
- Shows internal error details

#### Production Environment
- Shows generic error messages
- Hides stack traces
- Logs errors without sensitive data
- Hides internal error details

## Integration with Server

The error handler is automatically integrated with the HTTP server through the dependency injection system:

```go
// In internal/presentation/http/server.go
func (s *Server) setupErrorHandler() {
    factory := errorhandler.NewErrorHandlerFactory(s.config)
    errorHandler := factory.CreateErrorHandler()
    
    s.echo.HTTPErrorHandler = errorHandler.HandleError
}
```

## Testing

### Unit Tests

Run unit tests:

```bash
go test ./internal/presentation/http/errorhandler/... -v
```

### Integration Tests

The integration tests verify that the error handler works correctly with the Echo server:

```bash
go test ./internal/presentation/http/errorhandler/... -v -run TestErrorHandlerIntegration
```

## Best Practices

### 1. Error Response Format
- Always use consistent error response format
- Include error code, message, and status code
- Add stack traces only in development
- Never expose sensitive information in production

### 2. Error Logging
- Log all errors with context (request details)
- Include stack traces for debugging
- Use structured logging
- Don't log sensitive information

### 3. Security
- Sanitize error messages in production
- Don't expose internal error details
- Use appropriate HTTP status codes
- Handle errors gracefully

### 4. Performance
- Keep error handling lightweight
- Avoid expensive operations in error handlers
- Use efficient logging
- Cache error responses when appropriate

## Examples

### Custom Error Handler

```go
// Create custom error handler for specific use case
config := &errorhandler.Config{
    Environment: "development",
    ShowDetails: true,
    LogErrors:   true,
}

handler := errorhandler.NewErrorHandler(config)

// Use in specific route group
protected := e.Group("/api/v1")
protected.Use(errorhandler.ErrorHandlingMiddleware(handler))
```

### Error Handler with Custom Configuration

```go
// Create error handler with custom configuration
factory := errorhandler.NewErrorHandlerFactory(appConfig)
handler := factory.CreateErrorHandlerWithConfig(&errorhandler.Config{
    Environment: "staging",
    ShowDetails: false,
    LogErrors:   true,
})
```

## Future Extensions

The error handler system is designed to be easily extensible:

1. **Custom Error Types**: Add support for domain-specific errors
2. **Error Metrics**: Collect error metrics for monitoring
3. **Error Reporting**: Integrate with error reporting services
4. **Error Caching**: Cache common error responses
5. **Error Localization**: Support for localized error messages
6. **Error Correlation**: Add correlation IDs for error tracking 