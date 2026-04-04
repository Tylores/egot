# Testing Quick Reference

## Most Common Commands

```bash
# During development - fast feedback
make test-unit

# Before committing
make test-verbose

# Coverage check
make test-coverage

# All tests
make test
```

## Testing Patterns

### Basic Test
```go
func TestMyFeature(t *testing.T) {
    result := mypackage.MyFunction()
    assert.Equal(t, expected, result)
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

### Table-Driven
```go
func TestCases(t *testing.T) {
    tests := []struct{ name, input, want string }{
        {"case1", "a", "A"},
        {"case2", "b", "B"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, mypackage.Process(tt.input))
        })
    }
}
```

## Test Locations

| What | Where |
|------|-------|
| Unit tests | Next to `.go` file as `*_test.go` |
| Integration tests | `test/integration/` with `_integration_test.go` |
| Shared helpers | `test/testhelpers/` |
| Generated mocks | `test/mocks/` |

## Generate Mocks

```bash
./scripts/generate_mocks.sh internal/YourService/repository.go Repository
```

## Assertions (Testify)

```go
assert.Equal(t, expected, actual)
assert.NotEqual(t, unexpected, actual)
assert.Nil(t, obj)
assert.NotNil(t, obj)
assert.True(t, condition)
assert.False(t, condition)
assert.Error(t, err)
assert.NoError(t, err)
assert.Contains(t, str, substring)
assert.Greater(t, a, b)
```

## All Make Targets

| Command | Purpose |
|---------|---------|
| `make test` | All tests |
| `make test-unit` | Unit tests only (fast) |
| `make test-integration` | Integration tests |
| `make test-verbose` | Verbose + race detector |
| `make test-coverage` | Generate coverage report |
| `make test-short` | Short tests only |
| `make test-race` | Race detector |
| `make test-package PKG=./internal/Bill` | Single package |
| `make test-service SERVICE=core` | Single service |

## Key Files

- **Makefile** - Test commands
- **TESTING_GUIDE.md** - Full tutorial
- **TESTING_WORKFLOW.md** - Development workflow
- **test/testhelpers/assertions.go** - Helper functions
- **internal/Bill/handler/handler_test.go** - Example

## Test Helpers

```go
h := testhelpers.New(t)
h.NoError(err)
h.Equal(expected, actual)
h.True(condition)
h.Contains(str, substring)
```

## Skip Tests

```go
if testing.Short() {
    t.Skip("expensive test")
}
```

## Tips

- Run `make test-unit` during dev for fast feedback
- Run `make test-verbose` before commits
- Generate mocks with `./scripts/generate_mocks.sh`
- Use table-driven tests for multiple cases
- Organize tests in `*_test.go` files
- Reference `internal/Bill/handler/handler_test.go` for patterns

---

**Need more detail?** Read `TESTING_GUIDE.md` or `TESTING_WORKFLOW.md`
