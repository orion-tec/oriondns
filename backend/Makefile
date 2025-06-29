.PHONY: test test-unit test-integration test-db test-coverage clean build help

BINARY_NAME=oriondns
HTTPSERVER_BINARY=httpserver

# Test configuration
TEST_DB_URL ?= postgres://postgres@localhost:5432/oriondns_test?sslmode=disable
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build both DNS and HTTP server binaries
	@echo "Building DNS server..."
	@go build -o $(BINARY_NAME) ./cmd/dnsserver/main.go
	@echo "Building HTTP server..."
	@go build -o $(HTTPSERVER_BINARY) ./cmd/httpserver/main.go
	@echo "Build complete"

clean: ## Clean build artifacts and test files
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(HTTPSERVER_BINARY)
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)
	@echo "Clean complete"

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests only (no database required)
	@echo "Running unit tests..."
	@SKIP_DB_TESTS=true go test -v ./internal/... ./server/...

test-integration: ## Run integration tests (requires database)
	@echo "Running integration tests..."
	@echo "Using test database: $(TEST_DB_URL)"
	@TEST_DATABASE_URL=$(TEST_DB_URL) go test -v -tags=integration ./...

test-db: ## Run database tests only
	@echo "Running database tests..."
	@echo "Using test database: $(TEST_DB_URL)"
	@TEST_DATABASE_URL=$(TEST_DB_URL) go test -v ./internal/*/

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@TEST_DATABASE_URL=$(TEST_DB_URL) go test -v -coverprofile=$(COVERAGE_OUT) ./...
	@go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@go tool cover -func=$(COVERAGE_OUT) | tail -1

test-coverage-text: ## Show coverage in terminal
	@echo "Running tests with coverage..."
	@TEST_DATABASE_URL=$(TEST_DB_URL) go test -v -coverprofile=$(COVERAGE_OUT) ./...
	@go tool cover -func=$(COVERAGE_OUT)

setup-test-db: ## Create test database (requires PostgreSQL running)
	@echo "Setting up test database..."
	@psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS oriondns_test;"
	@psql -h localhost -U postgres -c "CREATE DATABASE oriondns_test;"
	@echo "Test database created"

run-migrations-test: ## Run migrations on test database
	@echo "Running migrations on test database..."
	@cd migrations && tern migrate --config tern-test.conf

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

lint: ## Run linter (if available)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running linter..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

dev-setup: deps setup-test-db ## Setup development environment
	@echo "Development environment setup complete"

ci-test: lint test-coverage ## Run CI tests (linting, vetting, and coverage)

# Docker targets
docker-build: ## Build Docker images
	@echo "Building Docker images..."
	@docker build -f Dockerfile-web-backend -t oriondns-backend .
	@docker build -f Dockerfile-frontend -t oriondns-frontend ../frontend/

docker-test: ## Run tests in Docker
	@echo "Running tests in Docker..."
	@docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

.DEFAULT_GOAL := help