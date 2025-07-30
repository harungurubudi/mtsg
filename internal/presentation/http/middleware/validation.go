package middleware

import (
	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/labstack/echo/v4"
)

// ValidationMiddleware provides request validation
func ValidationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: Implement request validation logic
			// This could include:
			// - Content-Type validation
			// - Request size limits
			// - Required headers
			// - Query parameter validation
			// - Body validation (if applicable)

			// For now, just pass through to next handler
			return next(c)
		}
	}
}

// ContentTypeMiddleware validates the Content-Type header
func ContentTypeMiddleware(allowedTypes ...string) echo.MiddlewareFunc {
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"application/json"}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")

			// Skip validation for GET requests
			if c.Request().Method == "GET" {
				return next(c)
			}

			// Check if content type is allowed
			for _, allowedType := range allowedTypes {
				if contentType == allowedType {
					return next(c)
				}
			}

			return pkgerror.NewValidationError("unsupported content type")
		}
	}
}
