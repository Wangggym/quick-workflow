#!/bin/bash

# Test script for quick-workflow
# This script runs comprehensive tests

set -e

cd "$(dirname "$0")/.."

echo "🧪 Running Quick Workflow Tests"
echo "================================"
echo ""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Track test results
FAILED=0

run_test() {
    local test_name=$1
    local test_cmd=$2
    
    echo -n "Running $test_name... "
    
    if eval "$test_cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# 1. Check Go installation
run_test "Go installation" "go version"

# 2. Check dependencies
run_test "Dependencies" "go mod download && go mod verify"

# 3. Run go vet
echo -n "Running go vet... "
if go vet ./... 2>&1 | grep -v "no Go files" > /tmp/vet.log; then
    if [ -s /tmp/vet.log ]; then
        echo -e "${RED}✗${NC}"
        cat /tmp/vet.log
        FAILED=$((FAILED + 1))
    else
        echo -e "${GREEN}✓${NC}"
    fi
else
    echo -e "${GREEN}✓${NC}"
fi

# 4. Run tests with coverage
echo "Running unit tests with coverage..."
if go test -v -race -coverprofile=coverage.out ./... 2>&1 | tee /tmp/test.log; then
    echo -e "${GREEN}✓ All tests passed${NC}"
    
    # Show coverage
    echo ""
    echo "Test Coverage:"
    go tool cover -func=coverage.out | tail -n 1
else
    echo -e "${RED}✗ Some tests failed${NC}"
    FAILED=$((FAILED + 1))
fi

# 5. Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
    echo ""
    echo "Running linters..."
    if golangci-lint run ./...; then
        echo -e "${GREEN}✓ Linters passed${NC}"
    else
        echo -e "${RED}✗ Linter issues found${NC}"
        FAILED=$((FAILED + 1))
    fi
else
    echo -e "${YELLOW}⚠ golangci-lint not installed, skipping linters${NC}"
fi

# 6. Build test
echo ""
echo -n "Testing build... "
if go build -o /tmp/qkflow-test ./cmd/qkflow; then
    echo -e "${GREEN}✓${NC}"
    rm /tmp/qkflow-test
else
    echo -e "${RED}✗${NC}"
    FAILED=$((FAILED + 1))
fi

# Summary
echo ""
echo "================================"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED test(s) failed${NC}"
    exit 1
fi

