# Microservices Testing Guide

This guide explains how to write and run tests for the egot microservices project using Go, testify, and gomock.

## Quick Start

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run integration tests
make test-integration

# Run tests with coverage report
make test-coverage

# Show all available test commands
make help
```

## Testing Framework

The project uses:
- **testify/assert**: For readable assertions in tests
- **testify/require**: For assertions that halt execution on failure
- **gomock**: For mocking interfaces and generating mock objects
- **table-driven tests**: For clean, scalable test organization

## Project Structure

```
egot/
├── cmd/                    # Microservice entry points
│   ├── core/
│   ├── crawler/
│   ├── flowreservation/
│   ├── operator/
│   ├── rsps/
│   └── ...
├── internal/               # Business logic and domain models
│   ├── Bill/
│   ├── Messaging/
│   ├── TariffProfile/
│   ├── TimeOfUse/
│   └── ...
└── test/
    ├── testhelpers/        # Shared test utilities
    ├── mocks/              # Generated mocks
    └── integration/        # Integration tests
```

## Writing Unit Tests

### Basic Unit Test Example

```go
package mypackage_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyFunction(t *testing.T) {
	// Arrange
	input := "test"
	expected := "TEST"
	
	// Act
	result := mypackage.MyFunction(input)
	
	// Assert
	assert.Equal(t, expected, result)
}
```

### Table-Driven Tests

```go
func TestMyFunctionWithMultipleCases(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "lowercase", input: "test", expected: "TEST"},
		{name: "uppercase", input: "TEST", expected: "TEST"},
		{name: "mixed", input: "TeSt", expected: "TEST"},
	}
	
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := mypackage.MyFunction(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
```

### Using the AssertHelper

The project provides a convenient AssertHelper for cleaner assertions:

```go
func TestWithHelper(t *testing.T) {
	h := testhelpers.New(t)
	
	h.NoError(err)
	h.Equal(expected, actual)
	h.True(condition)
	h.Contains(str, substring)
}
```

## Creating Mocks with Gomock

### Generate a Mock

```bash
# Generate mock for an interface
mockgen -source=internal/MyPackage/interface.go -destination=test/mocks/mock_interface.go -package=mocks
```

### Using a Mock in Tests

```go
package mypackage_test

import (
	"testing"
	"github.com/golang/mock/gomock"
	"myproject/test/mocks"
)

func TestWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	// Create mock
	mockRepo := mocks.NewMockRepository(ctrl)
	
	// Set expectations
	mockRepo.EXPECT().
		GetUser(gomock.Any()).
		Return(&User{ID: 1}, nil).
		Times(1)
	
	// Use mock in test
	service := NewUserService(mockRepo)
	user, err := service.GetUser(1)
	
	// Verify
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
}
```

## Integration Tests

### Mark Tests as Integration

Use the `// +build integration` build tag:

```go
//go:build integration

package integration_test

import "testing"

func TestIntegrationWithDatabase(t *testing.T) {
	// This test runs when using `make test-integration`
}
```

Or name your test file `*_integration_test.go` and use:

```go
func TestIntegration_UserService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Integration test logic
}
```

### Run Integration Tests

```bash
make test-integration
```

## Best Practices

### 1. Arrange-Act-Assert Pattern
```go
func TestSomething(t *testing.T) {
	// Arrange: Set up test data
	data := setupTestData()
	
	// Act: Call the function
	result := functionUnderTest(data)
	
	// Assert: Verify the result
	assert.Equal(t, expected, result)
}
```

### 2. Use Subtests for Organization
```go
func TestUserService(t *testing.T) {
	t.Run("Create user", func(t *testing.T) { /* ... */ })
	t.Run("Update user", func(t *testing.T) { /* ... */ })
	t.Run("Delete user", func(t *testing.T) { /* ... */ })
}
```

### 3. Cleanup with defer
```go
func TestWithSetup(t *testing.T) {
	// Setup
	resource := setupResource()
	defer resource.Close() // Cleanup
	
	// Test
}
```

### 4. Use Constants for Test Data
```go
const (
	testUserID   = 1
	testUserName = "John Doe"
	testEmail    = "john@example.com"
)
```

### 5. Mock External Dependencies
```go
mockDB := mocks.NewMockDatabase(ctrl)
mockCache := mocks.NewMockCache(ctrl)
service := NewService(mockDB, mockCache)
```

## Running Tests

### All Tests
```bash
make test
```

### Unit Tests Only
```bash
make test-unit
```

### Integration Tests Only
```bash
make test-integration
```

### Specific Package
```bash
make test-package PKG=./internal/Bill
```

### Specific Service
```bash
make test-service SERVICE=core
```

### With Coverage Report
```bash
make test-coverage
```

### With Race Detector
```bash
make test-race
```

### Verbose Output
```bash
make test-verbose
```

## Coverage Reports

Generate coverage reports:

```bash
make test-coverage
# Opens coverage.html in your browser
```

View coverage for a specific package:
```bash
go test -cover ./internal/Bill/...
```

## Continuous Integration

The Makefile targets work well in CI/CD pipelines:

```bash
# Run all tests and generate coverage
make test-coverage

# Run with race detector
make test-race

# Run lint checks
make lint
```

## Common Patterns

### Testing Error Cases
```go
func TestErrorHandling(t *testing.T) {
	_, err := functionThatErrors()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected error message")
}
```

### Mocking HTTP Handlers
```go
// Use httptest package
import "net/http/httptest"

func TestHTTPHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	
	handler := http.HandlerFunc(myHandler)
	handler.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}
```

### Testing Database Operations
```go
func TestDatabaseQuery(t *testing.T) {
	// Use in-memory database or test database
	db := setupTestDB(t)
	defer db.Close()
	
	repo := NewRepository(db)
	result, err := repo.Query()
	assert.NoError(t, err)
}
```

## Troubleshooting

### Mock Generation Issues
```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Regenerate mocks
go generate ./...
```

### Tests Hanging
- Use `-timeout` flag: `go test -timeout 30s ./...`
- Check for goroutine leaks or deadlocks

### Coverage Not Accurate
- Ensure you're excluding vendor code: `go test -coverprofile=coverage.out ./...`
- Run all tests: `make test-coverage`

## Further Reading

- [Go Testing Package Docs](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [GoMock Documentation](https://github.com/golang/mock)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
