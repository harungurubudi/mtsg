package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/harungurubudi/mtsg/pkg/redis"
	"github.com/labstack/echo/v4"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerMinute int
	RequestsPerHour   int
	KeyPrefix         string
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RequestsPerMinute: 60,
		RequestsPerHour:   1000,
		KeyPrefix:         "rate_limit",
	}
}

// RateLimitMiddleware provides rate limiting functionality
func RateLimitMiddleware(redisAdapter redis.AdapterRepository, config *RateLimitConfig) echo.MiddlewareFunc {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip rate limiting if Redis adapter is not available
			if redisAdapter == nil {
				return next(c)
			}

			// Get client IP for rate limiting
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			// Check rate limit for minute
			if err := checkRateLimit(c.Request().Context(), redisAdapter, clientIP, "minute", config.RequestsPerMinute, config.KeyPrefix); err != nil {
				return err
			}

			// Check rate limit for hour
			if err := checkRateLimit(c.Request().Context(), redisAdapter, clientIP, "hour", config.RequestsPerHour, config.KeyPrefix); err != nil {
				return err
			}

			// Call next handler
			return next(c)
		}
	}
}

// checkRateLimit checks if the request is within the rate limit
func checkRateLimit(ctx context.Context, redisAdapter redis.AdapterRepository, clientIP, period string, limit int, keyPrefix string) error {
	// Create rate limit key
	key := fmt.Sprintf("%s:%s:%s:%s", keyPrefix, clientIP, period, time.Now().Format("2006-01-02-15-04"))

	// Get current count from Redis
	var countStr string
	err := redisAdapter.GetByKey(ctx, key, &countStr)
	if err != nil && err.Error() != "redis: nil" {
		// Log error but don't block request
		return nil
	}

	count := 0
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			// Log error but don't block request
			return nil
		}
	}

	// Check if limit exceeded
	if count >= limit {
		return pkgerror.NewHTTPError(429, "rate_limit_exceeded", fmt.Sprintf("rate limit exceeded: %d requests per %s", limit, period))
	}

	// Increment count
	newCount := count + 1
	err = redisAdapter.Set(ctx, key, strconv.Itoa(newCount), getExpiration(period))
	if err != nil {
		// Log error but don't block request
		return nil
	}

	// Set rate limit headers
	// Note: These would be set in the response, but we don't have access to response writer here
	// They should be set in the response handler or another middleware

	return nil
}

// getExpiration returns the expiration time for the given period
func getExpiration(period string) time.Duration {
	switch period {
	case "minute":
		return time.Minute
	case "hour":
		return time.Hour
	default:
		return time.Minute
	}
}
