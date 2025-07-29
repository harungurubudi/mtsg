package di

import (
	redisclient "github.com/go-redis/redis/v8"
	"github.com/harungurubudi/mtsg/internal/presentation/http"
	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/internal/repository"
	"github.com/harungurubudi/mtsg/internal/usecase"
	"github.com/harungurubudi/mtsg/pkg/redis"
	"github.com/harungurubudi/mtsg/pkg/token"
)

// Configuration Providers

// ProvideRedisClient provides a Redis client instance
func ProvideRedisClient() *redisclient.Client {
	return redisclient.NewClient(&redisclient.Options{
		Addr:     "localhost:6379", // Default Redis address
		Password: "",               // No password by default
		DB:       0,                // Default database
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
func ProvideHandlers() *handler.Handlers {
	return handler.NewHandlers()
}

// Server Providers

// ProvideServerConfig provides server configuration
func ProvideServerConfig() *http.Config {
	return http.NewConfig()
}

// ProvideHTTPServer provides HTTP server instance
func ProvideHTTPServer(
	handlers *handler.Handlers,
	config *http.Config,
) *http.Server {
	return http.NewServer(handlers, config)
}
