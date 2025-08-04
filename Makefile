# MTSG Makefile
# Multi-Tenant SaaS Gateway

# Variables
APP_NAME := mtsg
MAIN_FILE := cmd/main.go
BUILD_DIR := build
BINARY_NAME := $(APP_NAME)
GO := go
WIRE := /home/harun/go/bin/wire
SWAG := /home/harun/go/bin/swag

# Build flags
LDFLAGS := -ldflags="-s -w"
BUILD_FLAGS := -tags=wireinject

# Default target
.PHONY: help
help: ## Show this help message
	@echo "MTSG - Multi-Tenant SaaS Gateway"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development targets
.PHONY: dev
dev: generate ## Run the application in development mode
	@echo "🚀 Starting MTSG in development mode..."
	$(GO) run $(MAIN_FILE)

.PHONY: dev-background
dev-background: generate ## Run the application in background
	@echo "🚀 Starting MTSG in background..."
	$(GO) run $(MAIN_FILE) &

# Build targets
.PHONY: build
build: generate ## Build the application
	@echo "🔨 Building MTSG..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-all
build-all: generate ## Build for all platforms
	@echo "🔨 Building MTSG for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "✅ Multi-platform build completed"

# Generate targets
.PHONY: generate
generate: wire swagger ## Generate all code (wire, swagger, etc.)

.PHONY: wire
wire: ## Generate Wire dependency injection code
	@echo "🔧 Generating Wire code..."
	$(WIRE) ./internal/di
	@echo "✅ Wire code generated"

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "📚 Generating Swagger documentation..."
	$(SWAG) init -g $(MAIN_FILE) -o internal/presentation/http/docs
	@echo "✅ Swagger documentation generated"

# Test targets
.PHONY: test
test: ## Run all tests
	@echo "🧪 Running tests..."
	$(GO) test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: test-short
test-short: ## Run short tests
	@echo "🧪 Running short tests..."
	$(GO) test -v -short ./...

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

.PHONY: clean-swagger
clean-swagger: ## Clean swagger generated files
	@echo "🧹 Cleaning swagger files..."
	rm -rf internal/presentation/http/docs
	@echo "✅ Swagger files cleaned"

# Lint and format targets
.PHONY: lint
lint: ## Run linter
	@echo "🔍 Running linter..."
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	@echo "🎨 Formatting code..."
	$(GO) fmt ./...

.PHONY: fmt-check
fmt-check: ## Check code formatting
	@echo "🔍 Checking code formatting..."
	@if [ "$$($(GO) fmt ./... | wc -l)" -gt 0 ]; then \
		echo "❌ Code is not formatted. Run 'make fmt' to fix."; \
		exit 1; \
	else \
		echo "✅ Code is properly formatted"; \
	fi

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	$(GO) mod download
	@echo "✅ Dependencies downloaded"

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "📦 Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "✅ Dependencies updated"

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):latest .
	@echo "✅ Docker image built"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):latest

# Utility targets
.PHONY: kill-port-8080
kill-port-8080: ## Kill process using port 8080
	@echo "🔫 Killing process on port 8080..."
	@lsof -ti:8080 | xargs -r kill -9
	@echo "✅ Port 8080 freed"

.PHONY: check-deps
check-deps: ## Check if required tools are installed
	@echo "🔍 Checking required tools..."
	@command -v $(WIRE) >/dev/null 2>&1 || { echo "❌ Wire is not installed. Run: go install github.com/google/wire/cmd/wire@latest"; exit 1; }
	@command -v $(SWAG) >/dev/null 2>&1 || { echo "❌ Swag is not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; exit 1; }
	@echo "✅ All required tools are installed"

# Install targets
.PHONY: install-tools
install-tools: ## Install required development tools
	@echo "🔧 Installing development tools..."
	$(GO) install github.com/google/wire/cmd/wire@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Quick start
.PHONY: quick-start
quick-start: check-deps generate dev ## Quick start: check deps, generate code, and run dev server

# Production targets
.PHONY: run
run: build ## Build and run the application
	@echo "🚀 Running MTSG..."
	./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run-background
run-background: build ## Build and run the application in background
	@echo "🚀 Running MTSG in background..."
	./$(BUILD_DIR)/$(BINARY_NAME) &

# Development workflow
.PHONY: dev-setup
dev-setup: install-tools deps generate ## Complete development setup
	@echo "🎉 Development environment is ready!"
	@echo "Run 'make dev' to start the application"

# Default target
.DEFAULT_GOAL := help 