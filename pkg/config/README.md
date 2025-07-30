# Configuration Package

This package provides configuration management for the MTSG application using [Viper](https://github.com/spf13/viper).

## Features

- **Multi-source configuration**: Environment variables, config files, and defaults
- **Type-safe**: Automatic type conversion and validation
- **Environment-aware**: Automatic environment variable binding
- **Flexible**: Support for YAML, JSON, TOML, and other formats

## Usage

### Basic Usage

```go
import "github.com/harungurubudi/mtsg/pkg/config"

// Load configuration
cfg, err := config.Load()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Use configuration
fmt.Printf("Server port: %s\n", cfg.Server.Port)
fmt.Printf("Redis host: %s\n", cfg.Redis.Host)
```

### Configuration Structure

```go
type Config struct {
    Server ServerConfig `mapstructure:"server"`
    Redis  RedisConfig  `mapstructure:"redis"`
}
```

## Environment Variables

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MTSG_SERVER_PORT` | `8080` | HTTP server port |
| `MTSG_SERVER_READ_TIMEOUT` | `30s` | Server read timeout |
| `MTSG_SERVER_WRITE_TIMEOUT` | `30s` | Server write timeout |
| `MTSG_SERVER_IDLE_TIMEOUT` | `60s` | Server idle timeout |
| `MTSG_SERVER_ENVIRONMENT` | `development` | Application environment |

### Redis Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MTSG_REDIS_HOST` | `localhost` | Redis server host |
| `MTSG_REDIS_PORT` | `6379` | Redis server port |
| `MTSG_REDIS_PASSWORD` | `` | Redis password (empty by default) |
| `MTSG_REDIS_DB` | `0` | Redis database number |
| `MTSG_REDIS_POOL_SIZE` | `10` | Redis connection pool size |

## Configuration Files

The application supports configuration files in multiple formats. Create a `config.yaml` file in one of these locations:

- Current directory (`.`)
- `./config/` directory
- `/etc/mtsg/` directory

### Example config.yaml

```yaml
server:
  port: "8080"
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"
  environment: "development"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

## Configuration Priority

Configuration values are loaded in the following order (highest to lowest priority):

1. **Environment variables** (e.g., `MTSG_SERVER_PORT`)
2. **Configuration files** (e.g., `config.yaml`)
3. **Default values** (hardcoded in the application)

## Testing

The package includes comprehensive tests that cover:

- Default configuration loading
- Environment variable overrides
- Configuration validation
- Error handling for invalid values

Run tests with:

```bash
go test ./pkg/config
```

## Integration with Dependency Injection

To use the configuration with Wire dependency injection:

```go
// internal/di/providers.go
func ProvideConfig() (*config.Config, error) {
    return config.Load()
}

// internal/di/wire.go
var ConfigSet = wire.NewSet(
    ProvideConfig,
    // ... other providers
)
```

## Future Extensions

This configuration package is designed to be easily extensible. When you're ready to add more configuration sections:

1. **Database Configuration**: Add `DatabaseConfig` struct and related environment variables
2. **JWT Configuration**: Add `JWTConfig` struct for authentication settings
3. **Logging Configuration**: Add `LogConfig` struct for log level and format settings

Each new configuration section should follow the same patterns established in this package. 