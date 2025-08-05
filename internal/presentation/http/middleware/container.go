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

	// Debug: Log container initialization
	if container.Authentication == nil {
		panic("Container Authentication is nil - DI initialization failed")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContainerContextKey, container)
			c.Logger().Infof("Container set in context with key: %s", ContainerContextKey)
			return next(c)
		}
	}
}

// GetContainerFromContext extracts container from Echo context
func GetContainerFromContext(c echo.Context) (*provider.Container, bool) {
	container, ok := c.Get(ContainerContextKey).(provider.Container)

	// Debug: Log what we find in context
	if !ok {
		c.Logger().Errorf("Container not found in context with key: %s", ContainerContextKey)
		// Let's also check what's actually in the context
		allKeys := make([]string, 0)
		for _, key := range c.ParamNames() {
			allKeys = append(allKeys, key)
		}
		c.Logger().Errorf("Available keys in context: %v", allKeys)

		// Try to get the raw value to see what's there
		rawValue := c.Get(ContainerContextKey)
		c.Logger().Errorf("Raw value for key %s: %v (type: %T)", ContainerContextKey, rawValue, rawValue)
		return nil, false
	}

	return &container, true
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
