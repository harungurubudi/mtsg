//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/harungurubudi/mtsg/internal/usecase"
)

// HandlerSet groups all handler-related providers
var HandlerSet = wire.NewSet(
	// Configuration
	ProvideRedisClient,
	ProvideRedisAdapter,

	// Repositories
	ProvideUserRepository,
	ProvideTokenGenerator,

	// Use Cases
	ProvideAuthUseCase,
)

//go:generate wire
func InitializeAuthUseCase() usecase.Authentication {
	panic(wire.Build(HandlerSet))
}
