//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/harungurubudi/mtsg/internal/di/provider"
)

var UseCaseSet = wire.NewSet(
	provider.ConfigSet,
	// Configuration
	provider.ProvideRedisClient,
	provider.ProvideRedisAdapter,

	// Repositories
	provider.ProvideUserRepository,
	provider.ProvideTokenGenerator,

	// Use Cases
	provider.ProvideAuthUseCase,

	// Container
	provider.ProvideContainer,
)

//go:generate wire
func InitializeContainer() provider.Container {
	panic(wire.Build(UseCaseSet))
}
