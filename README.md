# Multi-Tenant Backend in Go

Backend skeleton for building **multi-tenant applications** in Go. It includes core building blocks like tenant isolation, authorization and authentication, tenant project scoping, and clean architecture

Designed and coded with help from [Cursor](https://cursor.sh) – an AI-assisted developer IDE.

---

## Features

- Multi-tenant structure with tenant scoping at every layer
- Clean Architecture with dependency injection using Google Wire
- JWT-based authentication with Redis token storage
- Swagger API documentation
- Comprehensive testing setup
- More features is coming soon

---

## Tech Stack

| Layer               | Tool                     |
|---------------------|--------------------------|
| Language            | Go 1.24.3                |
| Framework           | Echo                     |
| Database            | PostgreSQL               |
| Cache               | Redis                    |
| DI Framework        | Google Wire              |
| API Documentation   | Swagger                  |
| IDE                 | [Cursor](https://cursor.sh) (AI-assisted) |

---

## Docs

- [`docs/architecture.md`](./docs/architecture.md)
- [`docs/tenant.md`](./docs/tenant.md)
- [`docs/user-journey.md`](./docs/user-journey.md)
- [`docs/data-model.md`](./docs/data-model.md)
- [`internal/di/README.md`](./internal/di/README.md) - Dependency Injection Guide

## Quick Start

### Prerequisites

- Go 1.24.3 or later
- Redis server running locally (or configure remote Redis)

### One-Command Setup

```bash
# Complete development setup (installs tools, dependencies, generates code)
make dev-setup

# Start the application
make dev
```

### Manual Setup

1. **Install development tools:**
   ```bash
   make install-tools
   ```

2. **Setup configuration:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Download dependencies:**
   ```bash
   make deps
   ```

4. **Generate code:**
   ```bash
   make generate
   ```

5. **Run the application:**
   ```bash
   make dev
   ```

## Makefile Commands

The project includes a comprehensive Makefile for easy development workflow:

### 🚀 Development Commands

```bash
make dev              # Run in development mode
make dev-background   # Run in background
make quick-start      # Check deps, generate code, and run dev server
```

### 🔨 Build Commands

```bash
make build            # Build the application
make build-all        # Build for all platforms (Linux, macOS, Windows)
make run              # Build and run
make run-background   # Build and run in background
```

### 🔧 Code Generation

```bash
make generate         # Generate all code (wire, swagger)
make wire             # Generate Wire dependency injection code
make swagger          # Generate Swagger documentation
```

### 🧪 Testing

```bash
make test             # Run all tests
make test-coverage    # Run tests with coverage report
make test-short       # Run short tests
```

### 🧹 Cleanup

```bash
make clean            # Clean build artifacts
make clean-swagger    # Clean swagger generated files
make kill-port-8080   # Kill process using port 8080
```

### 📦 Dependency Management

```bash
make deps             # Download dependencies
make deps-update      # Update dependencies
```

### 🎨 Code Quality

```bash
make fmt              # Format code
make fmt-check        # Check code formatting
make lint             # Run linter
```

### 🐳 Docker

```bash
make docker-build     # Build Docker image
make docker-run       # Run Docker container
```

### 🔍 Utility Commands

```bash
make check-deps       # Check if required tools are installed
make help             # Show all available commands
```

## Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Modify the `.env` file with your configuration:
   ```bash
   # Server configuration
   MTSG_SERVER_PORT=8080
   MTSG_SERVER_ENVIRONMENT=development
   
   # Redis configuration
   MTSG_REDIS_HOST=localhost
   MTSG_REDIS_PORT=6379
   ```

## Available Endpoints

- `GET /ping` - Health check
- `GET /health` - Detailed health status
- `GET /swagger/index.html` - API documentation
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh token

## Development Workflow

1. **Start development:**
   ```bash
   make dev
   ```

2. **Make changes to code**

3. **Regenerate code if needed:**
   ```bash
   make generate
   ```

4. **Run tests:**
   ```bash
   make test
   ```

5. **Format code:**
   ```bash
   make fmt
   ```

## Troubleshooting

### Port 8080 Already in Use
```bash
make kill-port-8080
```

### Missing Tools
```bash
make install-tools
```

### Code Generation Issues
```bash
make clean
make generate
```

### Dependency Issues
```bash
make deps-update
```
