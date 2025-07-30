package handler

// Handlers holds all HTTP handlers
type Handlers struct {
	Ping *PingHandler
	// TODO: Add other handlers as they are implemented
	// Auth  *AuthHandler
	// User  *UserHandler
	// Tenant *TenantHandler
}

// NewHandlers creates a new Handlers instance
func NewHandlers() *Handlers {
	return &Handlers{
		Ping: NewPingHandler(),
		// TODO: Initialize other handlers
		// Auth:  NewAuthHandler(authUseCase),
		// User:  NewUserHandler(userUseCase),
		// Tenant: NewTenantHandler(tenantRepo),
	}
}
