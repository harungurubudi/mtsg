package errorhandler

import (
	"github.com/labstack/echo/v4"
)

// ErrorHandlingMiddleware provides error handling for middleware
func ErrorHandlingMiddleware(errorHandler *ErrorHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				errorHandler.HandleError(err, c)
				return nil // Error already handled
			}
			return nil
		}
	}
}

// ErrorHandlingMiddlewareWithConfig creates error handling middleware with configuration
func ErrorHandlingMiddlewareWithConfig(config *Config) echo.MiddlewareFunc {
	errorHandler := NewErrorHandler(config)
	return ErrorHandlingMiddleware(errorHandler)
}

// DevelopmentErrorHandlingMiddleware creates error handling middleware for development
func DevelopmentErrorHandlingMiddleware() echo.MiddlewareFunc {
	config := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   true,
	}
	return ErrorHandlingMiddlewareWithConfig(config)
}

// ProductionErrorHandlingMiddleware creates error handling middleware for production
func ProductionErrorHandlingMiddleware() echo.MiddlewareFunc {
	config := &Config{
		Environment: "production",
		ShowDetails: false,
		LogErrors:   true,
	}
	return ErrorHandlingMiddlewareWithConfig(config)
}
