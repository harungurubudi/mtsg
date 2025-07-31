package handler

import (
	"github.com/harungurubudi/mtsg/internal/usecase"
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
func NewHandlers(authUseCase usecase.Authentication) *Handlers {
	return &Handlers{
		Ping: NewPingHandler(),
		Auth: NewAuthHandler(authUseCase),
		// TODO: Initialize other handlers
		// User:  NewUserHandler(userUseCase),
		// Tenant: NewTenantHandler(tenantRepo),
	}
}
