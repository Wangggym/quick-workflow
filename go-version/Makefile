SHELL := /bin/bash

HIDE ?= @

.PHONY: gen fix check dev run test test-cov test-fast build build-all install local quick clean help

# Variables
BINARY_NAME := qkflow
BUILD_DIR := bin
GO := go
GOFLAGS := -v

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Initialize dependencies
gen:
	$(HIDE)echo "📚 Initializing project..."
	$(HIDE)$(GO) mod download
	$(HIDE)$(GO) mod tidy
	$(HIDE)$(GO) mod verify
	@echo "✅ Dependencies initialized and verified."

# Code formatting and fixes
fix:
	$(HIDE)echo "💅 Formatting code..."
	$(HIDE)$(GO) fmt ./...
	$(HIDE)if command -v golangci-lint > /dev/null; then \
		echo "🔧 Fixing linting issues..."; \
		golangci-lint run --fix ./...; \
	fi
	@echo "✅ Code formatted and fixed."

# Code checking (linting)
check:
	$(HIDE)echo "🔍 Running linters..."
	$(HIDE)if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Run: brew install golangci-lint"; \
		exit 1; \
	fi
	@echo "✅ Check completed."

# Development mode (quick run with optional args)
# Usage: make dev ARGS="pr create"
dev:
	$(HIDE)$(GO) run ./cmd/qkflow $(ARGS)

# Production build and run
run: build
	$(HIDE)$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	$(HIDE)echo "🧪 Running tests..."
	$(HIDE)$(GO) test -v -race ./...
	@echo "✅ Tests passed."

# Tests with coverage
test-cov:
	$(HIDE)echo "🧪 Running tests with coverage..."
	$(HIDE)$(GO) test -v -race -coverprofile=coverage.out ./...
	$(HIDE)$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated in coverage.html"

# Fast tests (parallel)
test-fast:
	$(HIDE)echo "🧪 Running fast tests..."
	$(HIDE)$(GO) test -short -race ./...
	@echo "✅ Fast tests completed."

# Build binary
build:
	$(HIDE)echo "🔨 Building $(BINARY_NAME)..."
	$(HIDE)mkdir -p $(BUILD_DIR)
	$(HIDE)$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/qkflow
	@echo "✅ Built: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all:
	$(HIDE)echo "🔨 Building for multiple platforms..."
	$(HIDE)mkdir -p $(BUILD_DIR)
	$(HIDE)GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/qkflow
	$(HIDE)GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/qkflow
	$(HIDE)GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/qkflow
	$(HIDE)GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/qkflow
	@echo "✅ Multi-platform builds completed."

# Install to GOPATH/bin
install: build
	$(HIDE)echo "📦 Installing $(BINARY_NAME)..."
	$(HIDE)cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

# Quick local build and install (one-command)
local: 
	@echo "⚡ Quick build and install..."
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/qkflow
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Built and installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"
	@echo "💡 Run 'qkflow version' to verify"

# Alias for local
quick: local

# Release helper
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Please set VERSION, e.g., make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "🚀 Creating release $(VERSION)..."
	@echo "1️⃣  Running tests..."
	@$(MAKE) test
	@echo "2️⃣  Running linters (optional)..."
	@$(MAKE) check || echo "⚠️  Linters not available, skipping..."
	@echo "3️⃣  Building all platforms..."
	@$(MAKE) build-all
	@echo ""
	@echo "✅ All checks passed! Ready to release:"
	@echo ""
	@echo "Run the following commands:"
	@echo "  git tag -a $(VERSION) -m 'Release $(VERSION)'"
	@echo "  git push origin $(VERSION)"
	@echo ""
	@echo "This will trigger GitHub Actions to build and publish the release."

# Cleanup
clean:
	-$(HIDE)rm -rf $(BUILD_DIR)
	-$(HIDE)rm -rf coverage.out coverage.html
	-$(HIDE)$(GO) clean
	@echo "✅ Cleaned up."

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make gen          - Initialize dependencies"
	@echo "  make fix          - Format and fix code"
	@echo "  make check        - Run linters"
	@echo "  make dev          - Development mode (quick run)"
	@echo "  make run          - Build and run binary"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run tests"
	@echo "  make test-cov     - Run tests with coverage"
	@echo "  make test-fast    - Fast tests (short mode)"
	@echo ""
	@echo "Building:"
	@echo "  make build        - Build binary"
	@echo "  make build-all    - Build for all platforms"
	@echo "  make install      - Install to GOPATH/bin"
	@echo "  make local        - ⚡ Quick build and install (one-step)"
	@echo "  make quick        - Alias for 'make local'"
	@echo ""
	@echo "Release:"
	@echo "  make release      - Create release (requires VERSION env)"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make help         - Show this help message"

