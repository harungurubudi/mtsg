# Multi-Tenant Backend in Go

Backend skeleton for building **multi-tenant applications** in Go. It includes core building blocks like tenant isolation, authorization and authentication, tenant project scoping, and clean architecture

Designed and coded with help from [Cursor](https://cursor.sh) – an AI-assisted developer IDE.

---

## Features

- Multi-tenant structure with tenant scoping at every layer
- More features is coming soon

---


## Tech Stack

| Layer               | Tool                     |
|---------------------|--------------------------|
| Language            | Go 1.24.3                |
| Framework           | Echo                     |
| Database            | PostgreSQL               |
| IDE                 | [Cursor](https://cursor.sh) (AI-assisted) |

---

## Docs

- [`docs/architecture.md`](./docs/architecture.md)
- [`docs/tenant.md`](./docs/tenant.md)
- [`docs/user-journey.md`](./docs/user-journey.md)
- [`docs/data-model.md`](./docs/data-model.md)

## Getting Started

### Prerequisites

- Go 1.24.3 or later
- Redis server running locally (or configure remote Redis)

### Configuration

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

### Running the Application

```bash
# Run the application
go run cmd/main.go

# Or build and run
go build -o bin/mtsg cmd/main.go
./bin/mtsg
```

The server will start on `http://localhost:8080` by default.

### Available Endpoints

- `GET /ping` - Health check
- `GET /health` - Detailed health status
- `GET /swagger/index.html` - API documentation
