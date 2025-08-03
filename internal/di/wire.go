//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/harungurubudi/mtsg/internal/di/provider"
	"github.com/harungurubudi/mtsg/internal/presentation/http"
	"github.com/harungurubudi/mtsg/internal/usecase"
)

// ConfigSet groups all configuration-related providers
var ConfigSet = wire.NewSet(
	provider.ProvideConfig,
)

// HandlerSet groups all handler-related providers
var HandlerSet = wire.NewSet(
	ConfigSet,
	// Configuration
	provider.ProvideRedisClient,
	provider.ProvideRedisAdapter,

	// Repositories
	provider.ProvideUserRepository,
	provider.ProvideTokenGenerator,

	// Use Cases
	provider.ProvideAuthUseCase,

	// Handlers
	provider.ProvideHandlers,
)

// MiddlewareSet groups all middleware-related providers
var MiddlewareSet = wire.NewSet(
	provider.ProvideMiddlewareFactory,
)

// ServerSet groups all server-related providers
var ServerSet = wire.NewSet(
	HandlerSet,
	MiddlewareSet,
	provider.ProvideHTTPServer,
)

//go:generate wire
func InitializeAuthUseCase() usecase.Authentication {
	panic(wire.Build(HandlerSet))
}

//go:generate wire
func InitializeServer() *http.Server {
	panic(wire.Build(ServerSet))
}
