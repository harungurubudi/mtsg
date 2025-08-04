package middleware

import (
	"github.com/harungurubudi/mtsg/internal/di"
	"github.com/harungurubudi/mtsg/internal/di/provider"
	http_error "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/labstack/echo/v4"
)

// ContainerContextKey is the key used to store DI container in Echo context
const ContainerContextKey = "di_container"

// ContainerMiddleware injects the DI container into Echo context
func ContainerMiddleware() echo.MiddlewareFunc {
	container := di.InitializeContainer()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContainerContextKey, container)
			return next(c)
		}
	}
}

// GetContainerFromContext extracts container from Echo context
func GetContainerFromContext(c echo.Context) (*provider.Container, bool) {
	container, ok := c.Get(ContainerContextKey).(*provider.Container)
	return container, ok
}

// RequireContainer middleware ensures that a valid container is present in context
func RequireContainer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			container, exists := GetContainerFromContext(c)
			if !exists || container == nil {
				return http_error.NewInternalServerError("DI container not available")
			}

			return next(c)
		}
	}
}
