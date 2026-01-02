# Makefile para Order Service

.PHONY: help build run test test-unit test-integration coverage clean

# Variables
APP_NAME=tc-golunch-order-service
BINARY_NAME=main
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")
COVERAGE_FILE=coverage.out

help: ## Show help
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o $(BINARY_NAME) cmd/api/main.go
	@echo "✅ Build completed!"

run: ## Run the application
	@echo "Starting $(APP_NAME) on port 8081..."
	go run cmd/api/main.go

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	go test -v -race -short ./internal/... ./tests/...
	@echo "✅ Unit tests completed!"

test-integration: ## Run integration tests  
	@echo "Running integration tests..."
	go test -v -race -tags=integration ./tests/...
	@echo "✅ Integration tests completed!"

coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	go test -coverprofile=$(COVERAGE_FILE) ./internal/... ./tests/...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	go tool cover -func=$(COVERAGE_FILE)
	@echo "✅ Coverage report generated: coverage.html"

coverage-check: ## Check if coverage meets minimum threshold (80%)
	@echo "Checking coverage threshold..."
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ "$${COVERAGE}" -lt 80 ]; then \
		echo "❌ Coverage is $${COVERAGE}%, minimum required is 80%"; \
		exit 1; \
	else \
		echo "✅ Coverage is $${COVERAGE}%, meets minimum requirement of 80%"; \
	fi

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./internal/...

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run --timeout=5m
	@echo "✅ Linting completed!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated!"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html
	go clean -cache
	@echo "✅ Cleaned!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .
	@echo "✅ Docker image built!"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8081:8080 --name $(APP_NAME) $(APP_NAME):latest

ci: deps fmt lint test coverage-check ## Run CI pipeline
	@echo "🚀 CI pipeline completed successfully!"

# Database setup for tests
setup-test-db: ## Setup test database
	@echo "Setting up test database..."
	docker run -d --name postgres-test \
		-e POSTGRES_USER=test_user \
		-e POSTGRES_PASSWORD=test_pass \
		-e POSTGRES_DB=golunch_test \
		-p 5433:5432 postgres:13
	@echo "✅ Test database is running on port 5433"

teardown-test-db: ## Teardown test database
	@echo "Stopping test database..."
	docker stop postgres-test || true
	docker rm postgres-test || true
	@echo "✅ Test database stopped"

# Integration test with real services
test-e2e: ## Run end-to-end tests
	@echo "Running end-to-end tests..."
	@echo "🔍 Testing Order Service endpoints..."
	curl -f http://localhost:8081/ping || (echo "❌ Order Service not responding" && exit 1)
	curl -f http://localhost:8081/customer/anonymous || (echo "❌ Anonymous endpoint not working" && exit 1)
	@echo "✅ E2E tests passed!"

# Development helpers
dev: ## Run in development mode with hot reload
	@echo "Starting development mode..."
	air -c .air.toml

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed!"