.PHONY: help build install uninstall test test-verbose test-cover clean run fmt vet deps tidy check

# Variables
BINARY_NAME=ntmux
GOPATH=$(shell go env GOPATH)
INSTALL_PATH=$(GOPATH)/bin/$(BINARY_NAME)

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✓ Build complete: ./$(BINARY_NAME)"

install: ## Install the CLI tool to GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	go install .
	@echo "✓ $(BINARY_NAME) installed successfully!"
	@echo "  Run '$(BINARY_NAME) --help' to get started"

uninstall: ## Uninstall the CLI tool from GOPATH/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)
	@echo "✓ $(BINARY_NAME) uninstalled"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	go test -v -count=1 ./...

test-cover: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

clean: ## Clean build artifacts and test cache
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@go clean -testcache
	@echo "✓ Clean complete"

run: ## Run the application (use ARGS="your args" to pass arguments)
	@go run . $(ARGS)

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ go vet passed"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

tidy: ## Tidy go.mod and go.sum
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

check: fmt vet test ## Run format, vet, and tests
