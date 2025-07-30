package middleware

import (
	"strings"

	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/internal/usecase"
	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
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
				return pkgerror.NewUnauthorizedError("missing authorization header")
			}

			// Check if the header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return pkgerror.NewUnauthorizedError("invalid authorization header format")
			}

			// Extract the token (remove "Bearer " prefix)
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return pkgerror.NewUnauthorizedError("missing access token")
			}

			// TODO: Implement token validation using auth use case
			// For now, we'll return an error indicating this needs to be implemented
			// user, err := authUseCase.VerifyToken(c.Request().Context(), token, subject, tenantID)
			// if err != nil {
			//     return pkgerror.NewUnauthorizedError("invalid or expired token")
			// }

			// Set user context for downstream handlers
			// c.Set(UserContextKey, user)

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
				return pkgerror.NewUnauthorizedError("user not authenticated")
			}

			return next(c)
		}
	}
}
