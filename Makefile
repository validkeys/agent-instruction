.PHONY: build test lint clean help

BINARY_NAME=agent-instruction
BUILD_DIR=build
GO=/usr/local/go/bin/go

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

lint: ## Run linters
	@echo "Running go vet..."
	@$(GO) vet ./...
	@echo "Running go fmt check..."
	@test -z $$($(GO) fmt ./...)

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean

all: lint test build ## Run lint, test, and build

.DEFAULT_GOAL := help
