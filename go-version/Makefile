SHELL := /bin/bash

HIDE ?= @

.PHONY: gen fix check dev run test test-cov build build-all install clean help

# Variables
BINARY_NAME := qkg
BUILD_DIR := bin
GO := go
GOFLAGS := -v

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Initialize dependencies
gen:
	$(HIDE)echo "📚 Downloading dependencies..."
	$(HIDE)$(GO) mod download
	$(HIDE)$(GO) mod tidy
	@echo "✅ Dependencies initialized."

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
	$(HIDE)$(GO) run ./cmd/qk $(ARGS)

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
	$(HIDE)$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/qk
	@echo "✅ Built: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all:
	$(HIDE)echo "🔨 Building for multiple platforms..."
	$(HIDE)mkdir -p $(BUILD_DIR)
	$(HIDE)GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/qk
	$(HIDE)GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/qk
	$(HIDE)GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/qk
	$(HIDE)GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/qk
	@echo "✅ Multi-platform builds completed."

# Install to GOPATH/bin
install: build
	$(HIDE)echo "📦 Installing $(BINARY_NAME)..."
	$(HIDE)cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

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
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts"

