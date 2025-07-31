package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/internal/usecase"
	http_error "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/labstack/echo/v4"
)

// UserContextKey is the key used to store user information in Echo context
const UserContextKey = "user"

// AuthMiddleware provides JWT-based authentication
func AuthMiddleware(authUseCase usecase.Authentication) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return http_error.NewUnauthorizedError("missing authorization header")
			}

			// Check if the header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return http_error.NewUnauthorizedError("invalid authorization header format")
			}

			// Extract the token (remove "Bearer " prefix)
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == "" {
				return http_error.NewUnauthorizedError("missing access token")
			}

			// Validate token using auth use case
			// For now, we'll use a default tenant ID since we don't have tenant context
			// In a real application, you might get tenant ID from subdomain, header, or other context
			user, err := authUseCase.VerifyToken(c.Request().Context(), token.Token(tokenStr), "access_token", tenant.TenantID{})
			if err != nil {
				return http_error.NewUnauthorizedError("invalid or expired token")
			}

			// Set user context for downstream handlers
			c.Set(UserContextKey, user)
			c.Set("user_id", uuid.UUID(user.ID).String())

			// Call next handler
			return next(c)
		}
	}
}

// GetUserFromContext extracts user information from Echo context
func GetUserFromContext(c echo.Context) (*user.User, bool) {
	user, ok := c.Get(UserContextKey).(*user.User)
	return user, ok
}

// RequireUser middleware ensures that a valid user is present in context
func RequireUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, exists := GetUserFromContext(c)
			if !exists || user == nil {
				return http_error.NewUnauthorizedError("user not authenticated")
			}

			return next(c)
		}
	}
}
