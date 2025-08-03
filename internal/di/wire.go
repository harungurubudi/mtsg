//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/harungurubudi/mtsg/internal/di/provider"
	"github.com/harungurubudi/mtsg/internal/usecase"
)

// ConfigSet groups all configuration-related providers
var ConfigSet = wire.NewSet(
	provider.ProvideConfig,
)

var AuthUseCaseSet = wire.NewSet(
	ConfigSet,
	// Configuration
	provider.ProvideRedisClient,
	provider.ProvideRedisAdapter,

	// Repositories
	provider.ProvideUserRepository,
	provider.ProvideTokenGenerator,

	// Use Cases
	provider.ProvideAuthUseCase,
)

//go:generate wire
func InitializeAuthUseCase() usecase.Authentication {
	panic(wire.Build(AuthUseCaseSet))
}
