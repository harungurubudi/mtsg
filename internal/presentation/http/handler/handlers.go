package handler

import (
	"github.com/labstack/echo/v4"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Ping *PingHandler
	Auth *AuthHandler
	// TODO: Add other handlers as they are implemented
	// User  *UserHandler
	// Tenant *TenantHandler
}

// NewHandlers creates a new Handlers instance
func NewHandlers() *Handlers {
	return &Handlers{
		Ping: NewPingHandler(),
		Auth: NewAuthHandler(),
		// TODO: Initialize other handlers
		// User:  NewUserHandler(userUseCase),
		// Tenant: NewTenantHandler(tenantRepo),
	}
}

func (h *Handlers) SetupRoutes(e *echo.Group) {
	// API routes (under /api/v1)
	h.Auth.SetupRoutes("auth", e)
}
