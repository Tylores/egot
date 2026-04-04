# Testing Setup for EGOT Microservices

## What's Included

This testing setup provides everything you need to write comprehensive unit and integration tests for your microservices:

### Core Components

✅ **testify** - Powerful assertions and mocking  
✅ **gomock** - Mock generation for interfaces  
✅ **Makefile targets** - Easy test commands  
✅ **Test helpers** - Reusable test utilities  
✅ **Example tests** - Starter templates  
✅ **Documentation** - Complete testing guides  

## Quick Start

```bash
# Run all tests
make test

# Run just unit tests (fast)
make test-unit

# Run with coverage report
make test-coverage

# See all options
make help
```

## Project Structure

```
egot/
├── Makefile                     # Test commands (make test, make test-unit, etc.)
├── TESTING_GUIDE.md             # Comprehensive testing tutorial
├── TESTING_WORKFLOW.md          # Recommended testing workflow
├── scripts/
│   └── generate_mocks.sh        # Helper script for mock generation
├── test/
│   ├── testhelpers/             # Shared test utilities
│   │   ├── assertions.go        # AssertHelper for cleaner assertions
│   │   └── testcase.go          # Table-driven test helpers
│   ├── mocks/                   # Generated mock files (will be created as needed)
│   └── integration/             # Integration tests
└── internal/
    └── Bill/
        └── handler/
            ├── handler.go
            └── handler_test.go  # Example unit test
```

## Key Files

| File | Purpose |
|------|---------|
| `Makefile` | Test commands and targets |
| `TESTING_GUIDE.md` | Complete testing tutorial with patterns |
| `TESTING_WORKFLOW.md` | Development workflow recommendations |
| `test/testhelpers/assertions.go` | Helper for readable assertions |
| `test/testhelpers/testcase.go` | Table-driven test utilities |
| `scripts/generate_mocks.sh` | Generate mocks from interfaces |
| `internal/Bill/handler/handler_test.go` | Example test file to reference |

## Common Tasks

### Write a Unit Test

1. Copy the pattern from `internal/Bill/handler/handler_test.go`
2. Create `<component>_test.go` next to your code
3. Use `testify/assert` for assertions
4. Run `make test-unit` for quick feedback

### Create a Mock

```bash
./scripts/generate_mocks.sh internal/YourService/repository.go Repository
```

### Run Tests for a Service

```bash
# Core service
make test-service SERVICE=core

# Bill service  
make test-service SERVICE=core

# Specific package
make test-package PKG=./internal/Messaging
```

### Generate Coverage Report

```bash
make test-coverage
# Opens coverage.html in your browser
```

## Make Targets Available

```
make test                 # Run all tests
make test-unit            # Fast unit tests only
make test-integration     # Integration tests
make test-verbose         # Verbose + race detector
make test-coverage        # Generate coverage report
make test-short           # Short tests only
make test-race            # Race detector only
make test-package PKG=... # Test specific package
make test-service SERVICE=...  # Test specific service
make generate-mocks       # Generate mocks info
make lint                 # Run linters
make help                 # Show all commands
```

## Writing Your First Test

### Simple Example

```go
package bill_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMyFunction(t *testing.T) {
	result := mypackage.MyFunction("input")
	assert.Equal(t, "expected", result)
}
```

### With Helper

```go
package bill_test

import (
	"testing"
	"github.com/Tylores/egot/test/testhelpers"
)

func TestMyFunction(t *testing.T) {
	h := testhelpers.New(t)
	
	result, err := mypackage.MyFunction("input")
	h.NoError(err)
	h.Equal("expected", result)
}
```

### Table-Driven Test

```go
func TestMultipleCases(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "case1", input: "a", expected: "A"},
		{name: "case2", input: "b", expected: "B"},
	}
	
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := mypackage.Process(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
```

## Testing Your Microservices

The project has 8 main microservices. Test each independently:

| Service | Command |
|---------|---------|
| core | `make test-service SERVICE=core` |
| crawler | `make test-service SERVICE=crawler` |
| flowreservation | `make test-service SERVICE=flowreservation` |
| operator | `make test-service SERVICE=operator` |
| rsps | `make test-service SERVICE=rsps` |
| client | `make test-service SERVICE=client` |
| scaffold-gen | `make test-service SERVICE=scaffold-gen` |
| wadl-extract | `make test-service SERVICE=wadl-extract` |

## Dependencies Added

```
github.com/stretchr/testify v1.11.1  # Assertions and mocking utilities
github.com/golang/mock v1.6.0        # Mock generation tool
```

## Next Steps

1. **Read the guides**
   - `TESTING_GUIDE.md` - Comprehensive patterns and examples
   - `TESTING_WORKFLOW.md` - Recommended development workflow

2. **Create your first test**
   - Copy `internal/Bill/handler/handler_test.go` as a template
   - Replace with your service's actual code
   - Run `make test-unit` to verify

3. **Generate mocks as needed**
   - Use `./scripts/generate_mocks.sh` when you need to mock dependencies
   - Example: `./scripts/generate_mocks.sh internal/Bill/repository/repository.go Repository`

4. **Set up coverage goals**
   - Aim for >75% on business logic
   - Use `make test-coverage` to track progress

5. **Integrate with CI/CD**
   - Run `make test` in your CI pipeline
   - Run `make test-coverage` for coverage reports
   - Run `make test-race` to detect data races

## Support & Resources

- **Testing Tutorial** → `TESTING_GUIDE.md`
- **Development Workflow** → `TESTING_WORKFLOW.md`
- **Example Tests** → `internal/Bill/handler/handler_test.go`
- **Mock Helper** → `./scripts/generate_mocks.sh`
- **Test Helpers** → `test/testhelpers/`

## Questions?

Refer to the comprehensive guides included in the project:

```bash
cat TESTING_GUIDE.md       # Complete tutorial
cat TESTING_WORKFLOW.md    # Recommended workflow
make help                  # Available commands
```

Happy testing! 🚀
