# Dependency Injection (DI) Setup

This directory contains the dependency injection configuration using Google Wire for compile-time dependency injection.

## Overview

We use **Google Wire** for compile-time dependency injection, which provides:
- **Type safety** - Compile-time verification of dependencies
- **Performance** - No runtime reflection overhead
- **Maintainability** - Clear dependency graph
- **Testability** - Easy to inject mocks
- **Clean Architecture** - Supports proper layer separation

## Directory Structure

```
internal/di/
├── provider/              # Provider functions organized by layer
│   ├── pkg.go            # Configuration and external service providers
│   ├── usecase.go        # Use case providers
│   └── presentation.go   # HTTP handlers, middleware, and server providers
├── wire.go               # Wire sets and initialization functions
├── wire_gen.go           # Generated DI code (auto-generated)
├── di_test.go            # Tests for DI setup
└── README.md             # This file
```

## Provider Organization

### Configuration & External Services (`provider/pkg.go`)
- `ProvideConfig()` - Application configuration
- `ProvideRedisClient()` - Redis client instance
- `ProvideRedisAdapter()` - Redis adapter wrapper
- `ProvideUserRepository()` - User data access
- `ProvideTokenGenerator()` - Token generation service

### Use Cases (`provider/usecase.go`)
- `ProvideAuthUseCase()` - Authentication business logic

### Presentation Layer (`provider/presentation.go`)
- `ProvideHandlers()` - HTTP handlers
- `ProvideMiddlewareFactory()` - Middleware factory
- `ProvideHTTPServer()` - HTTP server instance

## Wire Sets

### ConfigSet
Groups all configuration-related providers:
```go
var ConfigSet = wire.NewSet(
    provider.ProvideConfig,
)
```

### HandlerSet
Groups all handler-related providers:
```go
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
```

### MiddlewareSet
Groups all middleware-related providers:
```go
var MiddlewareSet = wire.NewSet(
    provider.ProvideMiddlewareFactory,
)
```

### ServerSet
Groups all server-related providers:
```go
var ServerSet = wire.NewSet(
    HandlerSet,
    MiddlewareSet,
    provider.ProvideHTTPServer,
)
```

## Usage

### Initialize Dependencies

```go
import "github.com/harungurubudi/mtsg/internal/di"

// Initialize authentication use case
authUseCase := di.InitializeAuthUseCase()

// Initialize complete HTTP server
server := di.InitializeServer()
```

### Add New Dependencies

1. **Add provider function** in the appropriate provider file:
   ```go
   // In provider/usecase.go
   func ProvideNewUseCase(repo repository.NewRepository) usecase.NewUseCase {
       return usecase.NewNewUseCase(repo)
   }
   ```

2. **Add to appropriate Wire Set** in `wire.go`:
   ```go
   var HandlerSet = wire.NewSet(
       // ... existing providers
       provider.ProvideNewUseCase,
   )
   ```

3. **Regenerate wire code**:
   ```bash
   /home/harun/go/bin/wire ./internal/di
   ```

## Provider Pattern

### Repository Providers
```go
// Provide concrete implementations
func ProvideUserRepository() repository.UserRepository {
    return repository.NewUserPersistence()
}

func ProvideTokenGenerator(redisAdapter redis.AdapterRepository) token.GeneratorRepository {
    return token.NewGenerator(redisAdapter, "your-secret-key-here")
}
```

### Use Case Providers
```go
// Inject dependencies into use cases
func ProvideAuthUseCase(
    userRepo repository.UserRepository,
    tokenGen token.GeneratorRepository,
) usecase.Authentication {
    return usecase.NewAuthentication(userRepo, tokenGen)
}
```

### Handler Providers
```go
// Inject use cases into handlers
func ProvideHandlers(authUseCase usecase.Authentication) *handler.Handlers {
    return handler.NewHandlers(authUseCase)
}
```

## Testing

### Run DI Tests
```bash
go test ./internal/di -v
```

### Mock Providers (Future Enhancement)
```go
// For testing, you can create mock providers
func ProvideMockUserRepository() repository.UserRepository {
    return &MockUserRepository{}
}

var MockHandlerSet = wire.NewSet(
    ProvideMockUserRepository,
    provider.ProvideAuthUseCase,
    provider.ProvideHandlers,
)
```

## Commands

### Generate Wire Code
```bash
# Using the installed wire binary
/home/harun/go/bin/wire ./internal/di

# Or add to PATH and use
go generate ./internal/di
```

### Test DI Setup
```bash
# Test compilation
go build ./internal/di

# Run tests
go test ./internal/di -v

# Test full application
go run cmd/main.go
```

## Architecture Principles

The DI setup follows Clean Architecture principles:

1. **Domain** - Business logic (no DI dependencies)
2. **Repository** - Data access (injected into use cases)
3. **Use Case** - Application logic (injected into handlers)
4. **Presentation** - HTTP handlers (injected with use cases)

All dependencies flow inward, maintaining proper layer separation.

## Naming Conventions

- **Provider Functions**: `Provide{Type}()` (e.g., `ProvideUserRepository()`)
- **Wire Sets**: `{Category}Set` (e.g., `HandlerSet`, `ServerSet`)
- **Initialization Functions**: `Initialize{Category}()` (e.g., `InitializeAuthUseCase()`)
- **Provider Files**: Organized by layer (`pkg.go`, `usecase.go`, `presentation.go`)

## Benefits

1. **Type Safety**: Compile-time verification of dependencies
2. **Performance**: No runtime reflection overhead
3. **Maintainability**: Clear dependency graph with modular organization
4. **Testability**: Easy to inject mocks for testing
5. **Scalability**: Easy to add new dependencies in appropriate files
6. **Clean Architecture**: Supports proper layer separation
7. **Modularity**: Providers organized by architectural layer

## Configuration

### Environment-based Configuration
```go
// Configuration is loaded from environment variables
func ProvideConfig() *config.Config {
    cfg, err := config.Load()
    if err != nil {
        panic(fmt.Sprintf("failed to load configuration: %v", err))
    }
    return cfg
}
```

Note: Wire generates clean, readable Go code that follows your project's patterns and conventions. 