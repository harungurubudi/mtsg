package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// CreateGlobalMiddleware returns all global middleware that should be applied to all routes
func CreateGlobalMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		ContainerMiddleware(), // Inject container first
		echoMiddleware.Recover(),
		echoMiddleware.Logger(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		}),
		echoMiddleware.RequestID(),
		echoMiddleware.Gzip(),
		HttpErrorHandler(),
	}
}

// CreateAuthMiddleware returns the authentication middleware
func CreateAuthMiddleware() echo.MiddlewareFunc {
	return AuthMiddleware()
}

// CreateRequireUserMiddleware returns middleware that ensures a valid user is present
func CreateRequireUserMiddleware() echo.MiddlewareFunc {
	return RequireUser()
}

// CreateRateLimitMiddleware returns rate limiting middleware
func CreateRateLimitMiddleware() echo.MiddlewareFunc {
	config := RateLimitConfig{
		RequestsPerMinute: 100,
		BurstSize:         10,
	}
	return RateLimitMiddleware(config)
}

// CreateValidationMiddleware returns request validation middleware
func CreateValidationMiddleware() echo.MiddlewareFunc {
	// TODO: Implement validation middleware
	// For now, return a no-op middleware
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}
