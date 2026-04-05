# Testing

## Quick Start

```bash
make test              # run all tests
make test-unit         # unit tests only (fast)
make test-integration  # integration tests
make test-coverage     # generate coverage report (opens coverage.html)
make test-verbose      # verbose output with race detector
make help              # show all targets
```

## Framework

- **testify/assert** — readable assertions
- **testify/require** — assertions that halt on failure
- **gomock** — interface mocking
- **table-driven tests** — scalable test organisation

Dependencies are in `go.mod`:
```
github.com/stretchr/testify v1.11.1
github.com/golang/mock v1.6.0
```

## Project Layout

```
egot/
├── Makefile
├── test/
│   ├── testhelpers/          # shared utilities
│   │   ├── assertions.go     # AssertHelper
│   │   └── testcase.go       # table-driven test helpers
│   ├── mocks/                # generated mocks
│   └── integration/          # integration tests
└── internal/
    └── Bill/handler/
        └── handler_test.go   # reference unit test
```

## Writing Tests

### Basic unit test

```go
package mypackage_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyFunction(t *testing.T) {
    result, err := mypackage.MyFunction("input")
    require.NoError(t, err)
    assert.Equal(t, "expected", result)
}
```

### Table-driven test

```go
func TestMultipleCases(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"lowercase", "test", "TEST"},
        {"uppercase", "TEST", "TEST"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, mypackage.Process(tt.input))
        })
    }
}
```

### Using the AssertHelper

```go
func TestWithHelper(t *testing.T) {
    h := testhelpers.New(t)
    result, err := mypackage.MyFunction("input")
    h.NoError(err)
    h.Equal("expected", result)
}
```

### Mocking HTTP handlers

```go
import "net/http/httptest"

func TestHTTPHandler(t *testing.T) {
    req, _ := http.NewRequest("GET", "/tp", nil)
    w := httptest.NewRecorder()
    handler := http.HandlerFunc(myHandler)
    handler.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Creating Mocks

```bash
./scripts/generate_mocks.sh internal/YourService/repository.go Repository
```

Using a mock in tests:

```go
func TestWithMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockRepository(ctrl)
    mockRepo.EXPECT().GetEntity(gomock.Any()).Return(entity, nil).Times(1)

    h := handler.NewHandler(mockRepo)
    // test handler...
}
```

## Integration Tests

Mark integration tests with a build tag or short-mode skip:

```go
//go:build integration
package integration_test
```

Or:

```go
func TestServiceIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    // ...
}
```

Run with:

```bash
make test-integration
```

## All Make Targets

| Command | Purpose |
|---------|---------|
| `make test` | All tests |
| `make test-unit` | Unit tests only |
| `make test-integration` | Integration tests |
| `make test-verbose` | Verbose + race detector |
| `make test-coverage` | Generate coverage HTML report |
| `make test-short` | Short tests only |
| `make test-race` | Race detector |
| `make test-package PKG=./internal/Bill` | Single package |
| `make test-service SERVICE=core` | Single service |
| `make generate-mocks` | Generate mocks info |
| `make lint` | Run linters |

## Development Workflow

**While writing new features:**
```bash
make test-unit       # fast feedback loop
```

**Before committing:**
```bash
make test-race       # check for race conditions
make test-coverage   # verify coverage
```

**Testing a specific service after changes:**
```bash
make test-service SERVICE=core
make test-package PKG=./internal/Messaging
```

## Coverage

Aim for:
- Critical paths / business logic: > 75%
- Error handling: > 70%

```bash
make test-coverage                                    # HTML report
go test -cover ./internal/Bill/...                   # per-package
go tool cover -func=coverage.out                     # detailed breakdown
```

## Best Practices

1. **Arrange-Act-Assert** — keep test sections clearly separated
2. **Use subtests** for grouped cases: `t.Run("case name", func(t *testing.T) { ... })`
3. **Cleanup with defer**: `defer resource.Close()`
4. **Mock external dependencies** — never hit real network or disk in unit tests
5. **Use `require` for fatal assertions** (stops test on failure), `assert` for non-fatal

## Troubleshooting

**Mock generation fails:**
```bash
go install github.com/golang/mock/mockgen@latest
go generate ./...
```

**Test hangs:**
```bash
go test -timeout 30s ./...
go test -race ./...
```

**Flaky tests in CI:**
```bash
go test -race -count=10 ./...
go test -race -shuffle=on ./...
```

## Service Test Commands

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

## Further Reading

- [quick-reference.md](quick-reference.md) — cheat sheet
- [routing.md](routing.md) — HTTP route testing
- [ci-cd.md](ci-cd.md) — CI/CD integration
- [internal/Bill/handler/handler_test.go](../../internal/Bill/handler/handler_test.go) — reference unit test
- [test/testhelpers/](../../test/testhelpers/) — shared test utilities
