package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/domain/tenant"
	"github.com/harungurubudi/mtsg/internal/domain/user"
	http_error "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/labstack/echo/v4"
)

// UserContextKey is the key used to store user information in Echo context
const UserContextKey = "user"

// ExtractTokenFromRequest extracts and validates the Bearer token from the request
func ExtractTokenFromRequest(c echo.Context) (token.Token, error) {
	// Extract token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", http_error.NewUnauthorizedError("missing authorization header")
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", http_error.NewUnauthorizedError("invalid authorization header format")
	}

	// Extract the token (remove "Bearer " prefix)
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		return "", http_error.NewUnauthorizedError("missing access token")
	}

	return token.Token(tokenStr), nil
}

// ValidateTokenAndGetUser validates the token and returns the user
func ValidateTokenAndGetUser(c echo.Context, tokenStr token.Token) (*user.User, error) {
	// Get container from context
	container, exists := GetContainerFromContext(c)
	if !exists {
		return nil, http_error.NewInternalServerError("DI container not available")
	}

	// Validate token using auth use case from container
	// For now, we'll use a default tenant ID since we don't have tenant context
	// In a real application, you might get tenant ID from subdomain, header, or other context
	user, err := container.Authentication.VerifyToken(c.Request().Context(), tokenStr, "access_token", tenant.TenantID{})
	if err != nil {
		return nil, http_error.NewUnauthorizedError("invalid or expired token")
	}

	return user, nil
}

// AuthMiddleware provides JWT-based authentication
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from request
			tokenStr, err := ExtractTokenFromRequest(c)
			if err != nil {
				return err
			}

			// Validate token and get user
			user, err := ValidateTokenAndGetUser(c, tokenStr)
			if err != nil {
				return err
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
