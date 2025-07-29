# Dependency Injection (DI) Setup

This directory contains the dependency injection configuration using Google Wire.

## Overview

We use **Google Wire** for compile-time dependency injection, which provides:
- **Type safety** - Compile-time verification of dependencies
- **Performance** - No runtime reflection overhead
- **Maintainability** - Clear dependency graph
- **Testability** - Easy to inject mocks

## Files

- `providers.go` - Provider functions for all dependencies
- `wire.go` - Wire sets and initialization functions (with build tags)
- `wire_gen.go` - Generated DI code (auto-generated)
- `di_test.go` - Tests for DI setup

## Current Dependencies

### Configuration
- `ProvideRedisClient()` - Redis client instance
- `ProvideRedisAdapter()` - Redis adapter wrapper

### Repositories
- `ProvideUserRepository()` - User data access
- `ProvideTokenGenerator()` - Token generation service

### Use Cases
- `ProvideAuthUseCase()` - Authentication business logic

## Usage

### Initialize Dependencies
```go
import "github.com/harungurubudi/mtsg/internal/di"

// Initialize authentication use case
authUseCase := di.InitializeAuthUseCase()
```

### Add New Dependencies

1. **Add provider function** in `providers.go`:
```go
func ProvideNewService() NewService {
    return NewService()
}
```

2. **Add to HandlerSet** in `wire.go`:
```go
var HandlerSet = wire.NewSet(
    // ... existing providers
    ProvideNewService,
)
```

3. **Regenerate wire code**:
```bash
/home/harun/go/bin/wire ./internal/di
```

## Testing

Run the DI tests:
```bash
go test ./internal/di -v
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
go run cmd/main.go
```

## Architecture

The DI setup follows Clean Architecture principles:
- **Domain** - Business logic (no DI dependencies)
- **Repository** - Data access (injected into use cases)
- **Use Case** - Application logic (injected into handlers)
- **Presentation** - HTTP handlers (injected with use cases)

All dependencies flow inward, maintaining proper layer separation. 