package middleware

import (
	"github.com/harungurubudi/mtsg/internal/usecase"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/harungurubudi/mtsg/pkg/redis"
	"github.com/labstack/echo/v4"
)

// MiddlewareFactory provides factory functions for middleware
type MiddlewareFactory struct {
	authUseCase  usecase.Authentication
	redisAdapter redis.AdapterRepository
	config       *config.Config
}

// NewMiddlewareFactory creates a new middleware factory
func NewMiddlewareFactory(
	authUseCase usecase.Authentication,
	redisAdapter redis.AdapterRepository,
	config *config.Config,
) *MiddlewareFactory {
	return &MiddlewareFactory{
		authUseCase:  authUseCase,
		redisAdapter: redisAdapter,
		config:       config,
	}
}

// GlobalMiddleware returns all global middleware
func (f *MiddlewareFactory) GlobalMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		f.RecoveryMiddleware(),
		f.LoggingMiddleware(),
		f.CORSMiddleware(),
	}
}

// AuthMiddleware returns authentication middleware
func (f *MiddlewareFactory) AuthMiddleware() echo.MiddlewareFunc {
	return AuthMiddleware(f.authUseCase)
}

// RateLimitMiddleware returns rate limiting middleware
func (f *MiddlewareFactory) RateLimitMiddleware() echo.MiddlewareFunc {
	return RateLimitMiddleware(f.redisAdapter, nil) // Use default config
}

// RateLimitMiddlewareWithConfig returns rate limiting middleware with custom config
func (f *MiddlewareFactory) RateLimitMiddlewareWithConfig(config *RateLimitConfig) echo.MiddlewareFunc {
	return RateLimitMiddleware(f.redisAdapter, config)
}

// ValidationMiddleware returns validation middleware
func (f *MiddlewareFactory) ValidationMiddleware() echo.MiddlewareFunc {
	return ValidationMiddleware()
}

// ContentTypeMiddleware returns content type validation middleware
func (f *MiddlewareFactory) ContentTypeMiddleware(allowedTypes ...string) echo.MiddlewareFunc {
	return ContentTypeMiddleware(allowedTypes...)
}

// RecoveryMiddleware returns recovery middleware
func (f *MiddlewareFactory) RecoveryMiddleware() echo.MiddlewareFunc {
	return RecoveryMiddleware()
}

// LoggingMiddleware returns logging middleware
func (f *MiddlewareFactory) LoggingMiddleware() echo.MiddlewareFunc {
	return LoggingMiddleware()
}

// CORSMiddleware returns CORS middleware
func (f *MiddlewareFactory) CORSMiddleware() echo.MiddlewareFunc {
	return CORSMiddleware()
}
