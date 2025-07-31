package di

import (
	"fmt"

	redisclient "github.com/go-redis/redis/v8"
	"github.com/harungurubudi/mtsg/internal/presentation/http"
	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/internal/repository"
	"github.com/harungurubudi/mtsg/internal/usecase"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/harungurubudi/mtsg/pkg/redis"
	"github.com/harungurubudi/mtsg/pkg/token"
)

// Configuration Providers

// ProvideConfig provides application configuration
func ProvideConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}

// ProvideRedisClient provides a Redis client instance
func ProvideRedisClient(cfg *config.Config) *redisclient.Client {
	return redisclient.NewClient(&redisclient.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
}

// ProvideRedisAdapter provides a Redis adapter instance
func ProvideRedisAdapter(client *redisclient.Client) redis.AdapterRepository {
	return redis.NewAdapter(client)
}

// Repository Providers

// ProvideUserRepository provides a concrete implementation of UserRepository
func ProvideUserRepository() repository.UserRepository {
	return repository.NewUserPersistence()
}

// ProvideTokenGenerator provides a concrete implementation of token.GeneratorRepository
func ProvideTokenGenerator(redisAdapter redis.AdapterRepository) token.GeneratorRepository {
	return token.NewGenerator(redisAdapter, "your-secret-key-here") // TODO: Use environment variable
}

// Use Case Providers

// ProvideAuthUseCase injects dependencies into Authentication usecase
func ProvideAuthUseCase(
	userRepo repository.UserRepository,
	tokenGen token.GeneratorRepository,
) usecase.Authentication {
	return usecase.NewAuthentication(userRepo, tokenGen)
}

// Handler Providers

// ProvideHandlers provides HTTP handlers
func ProvideHandlers(authUseCase usecase.Authentication) *handler.Handlers {
	return handler.NewHandlers(authUseCase)
}

// Server Providers

// ProvideHTTPServer provides HTTP server instance
func ProvideHTTPServer(
	handlers *handler.Handlers,
	config *config.Config,
) *http.Server {
	return http.NewServer(handlers, config)
}
