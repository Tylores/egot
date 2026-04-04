# EGOT Microservices Testing - Complete Setup Guide

## Overview

Your microservices project now has a complete, production-ready testing framework for all 8 services: **core**, **crawler**, **flowreservation**, **operator**, **rsps**, **client**, **scaffold-gen**, and **wadl-extract**.

## 📚 Documentation Index

### Start Here
- **[TESTING_SETUP_README.md](./TESTING_SETUP_README.md)** - Quick overview and getting started guide

### Detailed Guides
- **[TESTING_GUIDE.md](./TESTING_GUIDE.md)** - Comprehensive testing tutorial with patterns and examples
- **[TESTING_WORKFLOW.md](./TESTING_WORKFLOW.md)** - Recommended development workflow
- **[TESTING_QUICK_REFERENCE.md](./TESTING_QUICK_REFERENCE.md)** - Quick cheat sheet for commands and patterns
- **[CI_CD_INTEGRATION.md](./CI_CD_INTEGRATION.md)** - Integrating tests with GitHub Actions, GitLab CI, Jenkins, etc.

## 🚀 Quick Start

### Most Common Commands

```bash
# Fast feedback during development
make test-unit

# Complete check before committing
make test-verbose

# Generate coverage report
make test-coverage

# See all options
make help
```

### Your First Test

1. Copy the pattern from: `internal/Bill/handler/handler_test.go`
2. Create `<component>_test.go` in your package
3. Write your test using testify assertions
4. Run `make test-unit` for fast feedback

### Generate Mocks

```bash
./scripts/generate_mocks.sh internal/YourService/repository.go Repository
```

## 📦 What's Included

### Test Framework
- ✅ **testify v1.11.1** - Assertions and mocking utilities
- ✅ **gomock v1.6.0** - Mock code generation
- ✅ **Go 1.23.0** - Latest stable Go version

### Tools & Commands
- ✅ **Makefile** - 13 test commands (test, test-unit, test-coverage, etc.)
- ✅ **Mock Helper Script** - `scripts/generate_mocks.sh` for generating mocks
- ✅ **Test Helpers** - Reusable assertion and test utilities in `test/testhelpers/`

### Documentation (5 Guides)
- ✅ **Setup Guide** - Getting started
- ✅ **Testing Tutorial** - Complete patterns and examples
- ✅ **Development Workflow** - Recommended approach for development
- ✅ **Quick Reference** - Commands and patterns cheat sheet
- ✅ **CI/CD Guide** - Integration with pipelines

### Example
- ✅ **Bill Handler Test** - Reference template with 4 passing tests

## 🎯 Available Make Commands

| Command | Purpose | Use Case |
|---------|---------|----------|
| `make test` | Run all tests | CI/CD pipeline |
| `make test-unit` | Unit tests only (fast) | During development |
| `make test-integration` | Integration tests | Before commits |
| `make test-verbose` | All tests + race detector | Pre-commit check |
| `make test-coverage` | Generate coverage report | Weekly review |
| `make test-short` | Short tests only | Quick checks |
| `make test-race` | Race detector | Finding concurrency bugs |
| `make test-package PKG=...` | Test specific package | Targeted testing |
| `make test-service SERVICE=...` | Test specific service | Service development |
| `make generate-mocks` | Mock generation help | Reference |
| `make lint` | Run linters | Code quality |
| `make help` | Show all commands | Reference |

## 📂 Project Structure

```
egot/
├── Makefile                              # Test commands
├── TESTING_SETUP_README.md              # Start here
├── TESTING_GUIDE.md                     # Full tutorial
├── TESTING_WORKFLOW.md                  # Development workflow
├── TESTING_QUICK_REFERENCE.md           # Cheat sheet
├── CI_CD_INTEGRATION.md                 # Pipeline integration
│
├── go.mod / go.sum                      # Dependencies (updated)
├── scripts/
│   └── generate_mocks.sh                # Mock generation helper
│
├── test/
│   ├── testhelpers/
│   │   ├── assertions.go                # AssertHelper utility
│   │   └── testcase.go                  # Table-driven patterns
│   ├── mocks/                           # Generated mocks go here
│   └── integration/                     # Integration tests
│
└── internal/
    ├── Bill/
    │   └── handler/
    │       ├── handler.go
    │       └── handler_test.go          # Example test
    ├── Messaging/
    ├── TariffProfile/
    ├── TimeOfUse/
    └── ... (other services)
```

## 🔍 Testing Patterns

### Basic Unit Test
```go
func TestMyFunction(t *testing.T) {
    result := mypackage.MyFunction("input")
    assert.Equal(t, "expected", result)
}
```

### With Error Handling
```go
func TestWithError(t *testing.T) {
    result, err := mypackage.MyFunction()
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests
```go
func TestCases(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "a", "A"},
        {"case2", "b", "B"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, process(tt.input))
        })
    }
}
```

## 🔧 Configuration Files

### go.mod
Updated with test dependencies:
- `github.com/stretchr/testify v1.11.1`
- `github.com/golang/mock v1.6.0`

### Makefile
Provides convenient test commands and targets. Edit to customize for your workflow.

### scripts/generate_mocks.sh
Automated mock generation script. Usage:
```bash
./scripts/generate_mocks.sh <source_file> <InterfaceName>
```

## 🎓 Learning Path

1. **Day 1: Setup**
   - Read: `TESTING_SETUP_README.md`
   - Run: `make test-unit`
   - Verify: Tests pass ✓

2. **Day 2: Your First Test**
   - Reference: `internal/Bill/handler/handler_test.go`
   - Create: Your first `*_test.go` file
   - Run: `make test-unit`

3. **Day 3: Mocks & Advanced**
   - Read: `TESTING_GUIDE.md`
   - Generate: Mocks with `./scripts/generate_mocks.sh`
   - Write: Tests with mocks

4. **Day 4: Workflow**
   - Read: `TESTING_WORKFLOW.md`
   - Practice: Development workflow
   - Master: Fast feedback loop

5. **Day 5: CI/CD**
   - Read: `CI_CD_INTEGRATION.md`
   - Setup: GitHub Actions / GitLab CI / etc.
   - Integrate: Tests into pipeline

## 💡 Pro Tips

### Development Loop (Fastest)
```bash
# 1. Write code
# 2. Test immediately
make test-unit

# 3. Iterate quickly
# 4. Before commit
make test-verbose
```

### Code Review
```bash
# Generate coverage before PR
make test-coverage

# Check for race conditions
make test-race

# Run before pushing
make test-verbose
```

### Debugging
```bash
# Find flaky tests
go test -count=10 -race ./...

# See detailed output
make test-verbose

# Profile specific package
make test-package PKG=./internal/YourService
```

## 📊 Testing Your Services

Each of your 8 microservices can be tested individually:

```bash
make test-service SERVICE=core            # Core service
make test-service SERVICE=crawler         # Crawler service
make test-service SERVICE=flowreservation # Flow reservation
make test-service SERVICE=operator        # Operator service
make test-service SERVICE=rsps            # RSPS service
make test-service SERVICE=client          # Client utilities
make test-service SERVICE=scaffold-gen    # Scaffold generator
make test-service SERVICE=wadl-extract    # WADL extractor
```

## 🔐 Best Practices

1. **Run tests locally first**
   - `make test-unit` during development
   - `make test-verbose` before committing

2. **Name your tests clearly**
   - `TestValidInput` - what it tests
   - `TestErrorHandling` - behavior tested
   - Use subtests for organization

3. **Keep tests focused**
   - One assertion per test (when possible)
   - Mock external dependencies
   - Use table-driven tests for multiple cases

4. **Monitor coverage**
   - Target >75% for business logic
   - Run `make test-coverage` weekly
   - Focus on critical paths first

5. **Write before committing**
   - Each PR should include tests
   - Tests should pass locally
   - Coverage should meet standards

## 🆘 Troubleshooting

### Tests Won't Compile
```bash
go mod tidy
make test-unit
```

### Need to Skip a Test
```go
if testing.Short() {
    t.Skip("expensive operation")
}
```

### Test Hangs
```bash
go test -timeout 30s ./...
```

### Data Race Detected
```bash
make test-race
# Review the race detector output
```

## 📝 Next Steps

1. ✅ Framework is installed and verified
2. 📖 Read `TESTING_SETUP_README.md` (5 min)
3. 📝 Create your first test from the template (15 min)
4. 🔄 Run `make test-unit` and see it pass (1 min)
5. 🎯 Generate mocks for your dependencies (5 min)
6. 📊 Run `make test-coverage` to track progress (2 min)
7. 🚀 Integrate with your CI/CD pipeline (30 min)

## 📞 Reference

**Quick lookup:**
- Commands → `make help`
- Patterns → `TESTING_QUICK_REFERENCE.md`
- Tutorial → `TESTING_GUIDE.md`
- Workflow → `TESTING_WORKFLOW.md`
- CI/CD → `CI_CD_INTEGRATION.md`

**Example test:** `internal/Bill/handler/handler_test.go` (copy this pattern!)

## ✨ Summary

You now have:
- ✅ Production-ready test framework
- ✅ Easy-to-use make commands
- ✅ Comprehensive documentation
- ✅ Example tests to reference
- ✅ Mock generation helpers
- ✅ CI/CD integration guides

Everything is set up and ready to use. Start writing tests today!

---

**Last updated:** 2026-04-04  
**Framework version:** 1.0  
**Go version:** 1.23.0  
**Dependencies:** testify v1.11.1, gomock v1.6.0
