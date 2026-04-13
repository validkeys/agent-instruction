.PHONY: build test test-coverage lint clean install help

BINARY_NAME=agent-instruction
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin
GO=/usr/local/go/bin/go
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/agent-instruction
	@echo "Binary created at $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=$(COVERAGE_FILE) ./...
	@echo "Coverage report written to $(COVERAGE_FILE)"
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "HTML coverage report written to $(COVERAGE_HTML)"

lint: ## Run linters
	@echo "Running go vet..."
	@$(GO) vet ./...
	@echo "Running go fmt check..."
	@test -z $$($(GO) fmt ./...)

install: build ## Install binary to system
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed successfully. Run '$(BINARY_NAME)' to use."

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@$(GO) clean

all: lint test build ## Run lint, test, and build

.DEFAULT_GOAL := help
