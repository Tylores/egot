# Microservices Testing Workflow

This document outlines the recommended workflow for testing your egot microservices as you develop.

## Quick Command Reference

```bash
# Run tests during development
make test              # All tests
make test-unit         # Fast unit tests only
make test-integration  # Integration tests
make test-coverage     # With coverage report

# Detailed results
make test-verbose      # Verbose output with race detector
make test-race         # Detect data races

# For specific areas
make test-package PKG=./internal/Bill
make test-service SERVICE=core
```

## Development Workflow

### 1. While Writing New Features

```bash
# Quick feedback loop - run short tests only
make test-unit

# Before commit - run everything
make test-verbose
```

### 2. Before Committing

```bash
# Ensure no race conditions
make test-race

# Check coverage for critical paths
make test-coverage
```

### 3. Testing Specific Service

```bash
# After modifying a service
make test-service SERVICE=core

# After modifying a package
make test-package PKG=./internal/Messaging
```

## Writing Tests for Your Services

### Template for Unit Test

Create a file named `<component>_test.go` next to the component:

```go
package handler_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"your-module/internal/YourService/handler"
)

func TestYourHandler(t *testing.T) {
	// Arrange
	handler := handler.New()
	
	// Act
	result, err := handler.DoSomething()
	
	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
```

### Template for Integration Test

Create `*_integration_test.go` files in your test/integration directory:

```go
//go:build integration

package integration_test

import "testing"

func TestServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	
	// Set up multiple services
	// Test their interactions
}
```

### Generate Mocks for Dependencies

```bash
# Generate mock for an interface
./scripts/generate_mocks.sh internal/YourService/repository.go Repository

# Use in test
mockRepo := mocks.NewMockRepository(ctrl)
mockRepo.EXPECT().GetData().Return(data, nil)
```

## Test Organization by Service

Your 8 main microservices:

```
core            - Core platform logic
├── Start: make test-service SERVICE=core

crawler         - Data crawler service
├── Start: make test-service SERVICE=crawler

flowreservation - Flow reservation logic
├── Start: make test-service SERVICE=flowreservation

operator        - Operator service
├── Start: make test-service SERVICE=operator

rsps            - RSPS service
├── Start: make test-service SERVICE=rsps

client          - Client utilities
├── Start: make test-service SERVICE=client

scaffold-gen    - Scaffold generator
├── Start: make test-service SERVICE=scaffold-gen

wadl-extract    - WADL extractor
├── Start: make test-service SERVICE=wadl-extract
```

## CI/CD Integration

For your CI/CD pipeline:

```bash
#!/bin/bash

# Run all tests with race detector
make test-race

# Generate coverage
make test-coverage

# Check if coverage meets threshold
if grep -q "coverage: [0-8][0-9]\.[0-9]%" coverage.out; then
    echo "Coverage below 80%"
    exit 1
fi
```

## Code Coverage Strategy

### Recommended Minimums

- **Critical paths**: >80% coverage
- **Business logic**: >75% coverage
- **Error handling**: >70% coverage
- **Utilities**: >60% coverage

### View Coverage

```bash
# Generate HTML report
make test-coverage

# View for specific package
go test -cover ./internal/Bill/...

# Get detailed breakdown
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Troubleshooting Tests

### Test Hangs

```bash
# Run with timeout
go test -timeout 30s ./...

# Check for goroutine leaks
go test -race ./...
```

### Test Fails on CI but Passes Locally

```bash
# Run in same conditions as CI
go test -race -count=10 ./...

# With explicit seed
go test -race -shuffle=on ./...
```

### Need to Skip Tests

```go
// Skip single test
func TestExpensiveOperation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping expensive test")
	}
	// ...
}

// Skip build tag
//go:build !integration
package mypackage_test
```

## Performance Considerations

### Fast Iteration

```bash
# Fastest - only changed packages
go test ./internal/Bill/...

# Faster - short tests
make test-unit

# Normal - all tests
make test
```

### Parallel Execution

Tests run in parallel by default. Control with:

```bash
# Sequential (slower but more reliable for debugging)
go test -p=1 ./...

# Max parallel (faster on multi-core)
go test -parallel 16 ./...
```

## Monitoring Test Health

### Weekly Review

```bash
# Full coverage report
make test-coverage

# Check test times
go test -v -count=1 ./... | grep -E "(PASS|FAIL|---)"

# Look for flaky tests
go test -count=10 -race ./...
```

### Continuous Checks

Add to git pre-commit hook:

```bash
#!/bin/bash
echo "Running pre-commit tests..."
go test -short ./... || exit 1
go vet ./...
```

## Examples for Your Services

### Testing Bill Service

```bash
make test-service SERVICE=core     # Tests core package logic
go test -v ./internal/Bill/...     # Tests Bill service
make test-package PKG=./internal/Bill/handler
```

### Testing Messaging Service

```bash
make test-service SERVICE=core     # If in core
go test -v ./internal/Messaging/...
```

### Testing with Mocks

```bash
./scripts/generate_mocks.sh internal/Bill/repository/memory/repository.go Repository
go test -v ./internal/Bill/handler/...
```

## Next Steps

1. **Add tests gradually** - Don't try to test everything at once
2. **Start with critical paths** - Focus on business logic
3. **Use the templates** - Copy pattern from `internal/Bill/handler/handler_test.go`
4. **Generate mocks** - Use `./scripts/generate_mocks.sh` for dependencies
5. **Monitor coverage** - Run `make test-coverage` weekly
6. **Integrate into CI** - Add test targets to your CI/CD pipeline

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Docs](https://github.com/stretchr/testify)
- [GoMock Docs](https://github.com/golang/mock)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

## Support

For detailed information on testing patterns:

```bash
cat TESTING_GUIDE.md      # Comprehensive testing guide
./scripts/generate_mocks.sh # Generate mocks helper
make help                 # Show all make targets
```
