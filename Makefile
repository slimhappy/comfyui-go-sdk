.PHONY: help build test examples clean fmt lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all examples
	@echo "Building examples..."
	@cd examples/basic && go build -o ../../bin/basic
	@cd examples/websocket && go build -o ../../bin/websocket
	@cd examples/advanced && go build -o ../../bin/advanced
	@echo "Done! Binaries in ./bin/"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

examples: build ## Build and show example usage
	@echo ""
	@echo "Examples built successfully!"
	@echo ""
	@echo "Run examples:"
	@echo "  ./bin/basic      - Basic workflow submission"
	@echo "  ./bin/websocket  - WebSocket event monitoring"
	@echo "  ./bin/advanced   - Advanced features"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Done!"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Done!"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Done!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Done!"

install: ## Install the SDK
	@echo "Installing SDK..."
	@go install
	@echo "Done!"

.DEFAULT_GOAL := help
