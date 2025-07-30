//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/harungurubudi/mtsg/internal/presentation/http"
	"github.com/harungurubudi/mtsg/internal/usecase"
)

// ConfigSet groups all configuration-related providers
var ConfigSet = wire.NewSet(
	ProvideConfig,
)

// HandlerSet groups all handler-related providers
var HandlerSet = wire.NewSet(
	ConfigSet,
	// Configuration
	ProvideRedisClient,
	ProvideRedisAdapter,

	// Repositories
	ProvideUserRepository,
	ProvideTokenGenerator,

	// Use Cases
	ProvideAuthUseCase,

	// Handlers
	ProvideHandlers,
)

// ServerSet groups all server-related providers
var ServerSet = wire.NewSet(
	HandlerSet,
	ProvideHTTPServer,
)

//go:generate wire
func InitializeAuthUseCase() usecase.Authentication {
	panic(wire.Build(HandlerSet))
}

//go:generate wire
func InitializeServer() *http.Server {
	panic(wire.Build(ServerSet))
}
