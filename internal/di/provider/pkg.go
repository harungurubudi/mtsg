package provider

import (
	"fmt"

	redisclient "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/harungurubudi/mtsg/internal/repository"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/harungurubudi/mtsg/pkg/redis"
	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/jmoiron/sqlx"
)

// ProvideConfig provides application configuration
func ProvideConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}

// ConfigSet groups all configuration-related providers
var ConfigSet = wire.NewSet(
	ProvideConfig,
)

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

// ProvideSqlx provides a sqlx.DB instance using configuration values
func ProvideSqlx(cfg *config.Config) *sqlx.DB {
	var generateConnectionString = func() string {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
	}
	db, err := sqlx.Connect("postgres", generateConnectionString())
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	return db
}
