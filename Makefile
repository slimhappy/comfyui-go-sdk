.PHONY: help build test examples clean fmt lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all examples
	@echo "Building examples..."
	@mkdir -p bin
	@cd examples/basic && go build -o ../../bin/basic
	@cd examples/websocket && go build -o ../../bin/websocket
	@cd examples/advanced && go build -o ../../bin/advanced
	@cd examples/progress && go build -o ../../bin/progress
	@cd examples/execute_from_json && go build -o ../../bin/execute_from_json
	@cd examples/queue_management && go build -o ../../bin/queue_management
	@cd examples/history_operations && go build -o ../../bin/history_operations
	@cd examples/model_info && go build -o ../../bin/model_info
	@cd examples/image_operations && go build -o ../../bin/image_operations
	@cd examples/error_handling && go build -o ../../bin/error_handling
	@cd examples/integration_test && go build -o ../../bin/integration_test
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
	@echo "  ./bin/basic              - Basic workflow submission"
	@echo "  ./bin/websocket          - WebSocket event monitoring"
	@echo "  ./bin/advanced           - Advanced features"
	@echo "  ./bin/progress           - Real-time progress tracking with visual progress bar"
	@echo "  ./bin/execute_from_json  - Execute workflow from JSON file"
	@echo "  ./bin/queue_management   - Queue operations and management"
	@echo "  ./bin/history_operations - History retrieval and analysis"
	@echo "  ./bin/model_info         - Model and node information queries"
	@echo "  ./bin/image_operations   - Image upload and download"
	@echo "  ./bin/error_handling     - Error handling and retry logic"
	@echo "  ./bin/integration_test   - Comprehensive integration test suite"




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

build-execute-json: ## Build execute_from_json example only
	@echo "Building execute_from_json example..."
	@mkdir -p bin
	@cd examples/execute_from_json && go build -o ../../bin/execute_from_json
	@echo "Done! Binary: ./bin/execute_from_json"

run-queue: build ## Build and run queue management example
	@echo "Running queue management example..."
	@./bin/queue_management

run-history: build ## Build and run history operations example
	@echo "Running history operations example..."
	@./bin/history_operations

run-model-info: build ## Build and run model info example
	@echo "Running model info example..."
	@./bin/model_info

run-image-ops: build ## Build and run image operations example
	@echo "Running image operations example..."
	@./bin/image_operations

run-error-handling: build ## Build and run error handling example
	@echo "Running error handling example..."
	@./bin/error_handling

run-integration: build ## Build and run integration test
	@echo "Running integration test..."
	@./bin/integration_test

run-all-examples: build ## Build and run all new examples sequentially
	@echo "Running all new examples..."
	@echo ""
	@echo "=== Queue Management ==="
	@./bin/queue_management || true
	@echo ""
	@echo "=== History Operations ==="
	@./bin/history_operations || true
	@echo ""
	@echo "=== Model Info ==="
	@./bin/model_info || true
	@echo ""
	@echo "=== Image Operations ==="
	@./bin/image_operations || true
	@echo ""
	@echo "=== Error Handling ==="
	@./bin/error_handling || true
	@echo ""
	@echo "=== Integration Test ==="
	@./bin/integration_test || true

.DEFAULT_GOAL := help
