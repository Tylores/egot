# TEST_ROUTING_GUIDE.md - Comprehensive Microservices HTTP Testing Guide

## Overview

This guide documents the complete HTTP route testing framework for the EGOT microservices project. It covers route registration verification, HTTP method testing, and CI/CD integration across 23 microservices with 525+ routes.

## Table of Contents

1. [Framework Architecture](#framework-architecture)
2. [Phase 3: Route Registration Tests](#phase-3-route-registration-tests)
3. [Phase 4: HTTP Verb Tests](#phase-4-http-verb-tests)
4. [Phase 5: CI/CD Integration](#phase-5-cicd-integration)
5. [Phase 6: Advanced Testing](#phase-6-advanced-testing)
6. [Running Tests](#running-tests)
7. [Troubleshooting](#troubleshooting)
8. [Extending for New Services](#extending-for-new-services)

## Framework Architecture

### Components

```
test/routing/
├── validator.go           # Core route validation
├── mux_inspector.go       # HTTP mux introspection
├── helpers.go             # Test helper functions
├── wadl_extractor.go      # WADL file parsing
├── test_data_helpers.go   # Test data loading and utilities
└── README.md              # Framework documentation

test/testdata/
├── bill_routes.json       # Test expectations (18 services)
├── edev_routes.json
├── ...
└── upt_routes.json

tools/
├── generate_test_routes.go       # Generate test data from WADL
├── generate_all_service_tests.sh # Generate test templates
└── generate_route_tests.go       # Generate test files

cmd/<Service>/
├── <service>_route_registration_test.go
├── <service>_http_methods_test.go (Phase 4)
└── main.go
```

### Key Types

- **RouteExpectation**: Defines an expected route (method + path)
- **RouteValidator**: Verifies routes are registered
- **HTTPMethodValidator**: Confirms HTTP method support
- **PathParameterValidator**: Validates path parameters
- **MuxInspector**: Inspects HTTP mux for registered routes

## Phase 3: Route Registration Tests

### Purpose

Verify that all routes defined in WADL files are correctly registered in each microservice.

### Test File Structure

Each service has `<service>_route_registration_test.go` with three test functions:

#### 1. TestXXXRouteRegistration

```go
func TestBillRouteRegistration(t *testing.T) {
    // 1. Create test mux
    mux := http.NewServeMux()
    repo := memory.NewRepository()
    h := NewHandler(repo)
    registerBillRoutes(mux, h)
    
    // 2. Load expectations from test data
    expectations, err := routing.LoadTestDataByServiceName("Bill")
    if err != nil {
        t.Skip(err)
    }
    
    // 3. Register expected routes with helper
    helper := routing.NewTestHelper(mux)
    helper.RegisterExpectedRoutes(expectations)
    
    // 4. Verify all routes exist
    all, missing := helper.AssertAllRoutesExist(expectations)
    if !all {
        t.Errorf("Missing %d routes", len(missing))
        for _, route := range missing {
            t.Errorf("  - %s", route)
        }
    }
    
    // 5. Verify route count
    matches, registered, expected := helper.AssertRouteCount(expectations)
    if !matches {
        t.Errorf("Count mismatch: expected %d, got %d", expected, registered)
    }
}
```

**What it tests:**
- ✓ All WADL routes are registered
- ✓ Route count matches WADL resource count
- ✓ No duplicate routes
- ✓ Correct method/path combinations

**Expected results:**
- Bill: 70 routes
- EDevice: 265 routes
- PPY, TP, UPT: 45 routes each
- etc.

#### 2. TestXXXPathParameters

```go
func TestBillPathParameters(t *testing.T) {
    validator := routing.NewPathParameterValidator()
    
    // Define expected path patterns
    patterns := map[string][]string{
        "/bill/{id1}":              {"id1"},
        "/bill/{id1}/ca/{id2}":     {"id1", "id2"},
        "/bill/{id1}/ca/{id2}/bp/{id3}": {"id1", "id2", "id3"},
    }
    
    for path, expectedParams := range patterns {
        validator.AddPathPattern(path, expectedParams)
        
        // Extract parameters
        extracted := validator.ExtractParameterNames(path)
        
        // Verify they match
        if len(extracted) != len(expectedParams) {
            t.Errorf("Path %s: param count mismatch", path)
        }
        for i, expected := range expectedParams {
            if extracted[i] != expected {
                t.Errorf("Path %s param %d: expected '%s', got '%s'",
                    path, i, expected, extracted[i])
            }
        }
    }
}
```

**What it tests:**
- ✓ Path parameters correctly extracted
- ✓ Parameter names match WADL specification
- ✓ Parameter order is correct
- ✓ No missing or extra parameters

#### 3. TestXXXHTTPMethods

```go
func TestBillHTTPMethods(t *testing.T) {
    // Setup
    mux := http.NewServeMux()
    repo := memory.NewRepository()
    h := NewHandler(repo)
    registerBillRoutes(mux, h)
    
    inspector := routing.NewMuxInspector(mux)
    
    // Test sample routes
    tests := []struct {
        method string
        path   string
        expect bool
    }{
        {"GET", "/bill", true},
        {"POST", "/bill", true},
        {"DELETE", "/bill/{id1}", true},
        {"PATCH", "/bill", false}, // Unsupported
    }
    
    for _, test := range tests {
        status, _ := inspector.TestRequest(test.method, test.path)
        found := status != http.StatusNotFound
        if found != test.expect {
            t.Errorf("%s %s: expected found=%v", 
                test.method, test.path, test.expect)
        }
    }
}
```

**What it tests:**
- ✓ Supported HTTP methods work (GET, POST, PUT, DELETE)
- ✓ Methods match WADL definitions
- ✓ Unsupported methods return appropriate status

### Running Phase 3 Tests

```bash
# Test specific service
go test -v ./cmd/Bill/... -run TestBillRouteRegistration

# Test all route registration tests
go test -v ./cmd/... -run "RouteRegistration"

# Test with coverage
go test -cover ./cmd/...
```

## Phase 4: HTTP Verb Tests

### Purpose

Verify that all HTTP verbs defined in WADL are properly supported and routed.

### Implementation Pattern

Each service adds `<service>_http_verbs_test.go`:

```go
func TestBillHTTPVerbs(t *testing.T) {
    // Load test data
    expectations, err := routing.LoadTestDataByServiceName("Bill")
    if err != nil {
        t.Skip(err)
    }
    
    // Group by method
    byMethod := routing.GroupRoutesByMethod(expectations)
    
    // Create test mux
    mux := http.NewServeMux()
    repo := memory.NewRepository()
    h := NewHandler(repo)
    registerBillRoutes(mux, h)
    
    inspector := routing.NewMuxInspector(mux)
    
    // Verify each method
    for method, routes := range byMethod {
        t.Run(method, func(t *testing.T) {
            for _, route := range routes {
                matched := inspector.AssertRouteMatchesMethod(method, route.Path)
                if !matched {
                    t.Errorf("%s %s not matched", method, route.Path)
                }
            }
        })
    }
}
```

### Comprehensive HTTP Method Testing

Create subtests for each HTTP verb:

```bash
go test -v ./cmd/Bill/... -run "HTTPMethods/GET"
go test -v ./cmd/Bill/... -run "HTTPMethods/POST"
go test -v ./cmd/Bill/... -run "HTTPMethods/PUT"
go test -v ./cmd/Bill/... -run "HTTPMethods/DELETE"
```

### Expected Verb Distribution

| Service   | GET | POST | PUT | DELETE | HEAD | Total |
|-----------|-----|------|-----|--------|------|-------|
| Bill      | 14  | 14   | 14  | 14     | 14   | 70    |
| EDevice   | 53  | 53   | 53  | 53     | 53   | 265   |
| PPY       | 9   | 9    | 9   | 9      | 9    | 45    |
| DERP      | 8   | 8    | 8   | 8      | 8    | 40    |

## Phase 5: CI/CD Integration

### Makefile Integration

Add to Makefile:

```makefile
.PHONY: test-routes
test-routes:
	@echo "Running route registration tests..."
	go test -v ./cmd/... -run "RouteRegistration"

.PHONY: test-methods
test-methods:
	@echo "Running HTTP method tests..."
	go test -v ./cmd/... -run "HTTPMethod"

.PHONY: test-all-routes
test-all-routes: test-routes test-methods
	@echo "All route tests passed!"

test: test-all-routes
	@echo "Running all tests..."
	go test -v ./...
```

### GitHub Actions Integration

Add to `.github/workflows/test.yml`:

```yaml
- name: Run route registration tests
  run: make test-routes

- name: Run HTTP method tests
  run: make test-methods

- name: Verify route counts
  run: |
    BILL_COUNT=$(grep -c '"Method"' test/testdata/bill_routes.json)
    EDEV_COUNT=$(grep -c '"Method"' test/testdata/edev_routes.json)
    echo "Bill routes: $BILL_COUNT"
    echo "EDevice routes: $EDEV_COUNT"
```

### Parallel Test Execution

```bash
# Run tests in parallel (8 goroutines)
go test -parallel 8 ./cmd/...

# With verbose output
go test -v -parallel 8 ./cmd/... -run "RouteRegistration"
```

## Phase 6: Advanced Testing

### Test Coverage Reports

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./cmd/...

# View coverage in browser
go tool cover -html=coverage.out -o coverage.html
open coverage.html

# Summary
go tool cover -func=coverage.out | grep total
```

### Performance Testing

```go
func BenchmarkRouteMatch(b *testing.B) {
    expectations, _ := routing.LoadTestDataByServiceName("Bill")
    helper := routing.NewTestHelper(http.NewServeMux())
    helper.RegisterExpectedRoutes(expectations)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        helper.AssertAllRoutesExist(expectations)
    }
}
```

Run: `go test -bench=./cmd/Bill/...`

### Integration Testing Pattern

Load data and test with actual handlers:

```go
func TestBillIntegration(t *testing.T) {
    expectations, _ := routing.LoadTestDataByServiceName("Bill")
    
    // Create real mux and handlers
    mux := http.NewServeMux()
    repo := memory.NewRepository()
    h := NewHandler(repo)
    registerBillRoutes(mux, h)
    
    // Test each route with sample data
    for _, route := range expectations {
        path := strings.ReplaceAll(route.Path, "{id1}", "test123")
        path = strings.ReplaceAll(path, "{id2}", "test456")
        
        req := httptest.NewRequest(route.Method, path, nil)
        w := httptest.NewRecorder()
        
        mux.ServeHTTP(w, req)
        
        // Status should not be 404 (unless path template issue)
        if w.Code == http.StatusNotFound {
            t.Errorf("%s %s returned 404", route.Method, route.Path)
        }
    }
}
```

## Running Tests

### Quick Start

```bash
# Run all route tests
make test-routes

# Run specific service
go test -v ./cmd/Bill/...

# Run with coverage
go test -cover ./cmd/...
```

### Comprehensive Testing

```bash
# All phases
go test -v ./cmd/... ./test/routing/

# Routes only
go test -v ./cmd/... -run "Route"

# HTTP methods only  
go test -v ./cmd/... -run "HTTPMethod"

# Path parameters only
go test -v ./cmd/... -run "PathParameter"
```

### Test Filtering

```bash
# Single service
go test -v ./cmd/Bill/... -run TestBillRouteRegistration

# Multiple services
go test -v ./cmd/Bill/... ./cmd/EDevice/... -run Route

# Specific test patterns
go test -v ./cmd/... -run "Route.*Registration"
```

## Troubleshooting

### Common Issues

#### 1. Routes Not Found

**Problem**: TestXXXRouteRegistration reports missing routes

**Diagnosis**:
```bash
# Check test data loads correctly
go run -run LoadTestDataByServiceName test/routing/test_data_helpers.go

# Verify WADL file exists
ls -la wadl/<service>.wadl

# Inspect test data
jq '.routes | length' test/testdata/<service>_routes.json
```

**Solution**:
1. Verify WADL files are present and valid
2. Regenerate test data: `go run ./tools/generate_test_routes.go`
3. Check registerXXXRoutes function matches WADL routes

#### 2. Path Parameter Issues

**Problem**: Path parameters not extracted correctly

**Example Error**: 
```
Path /bill/{id1}: param count mismatch - expected 1, got 0
```

**Solution**:
1. Verify parameter names in WADL match path syntax
2. Check WADL samplePath attribute format
3. Ensure parameter names follow {name} syntax

#### 3. HTTP Method Mismatches

**Problem**: Some HTTP methods report not found

**Causes**:
- Handler not registered for method/path
- Route template mismatch
- Parameter syntax incorrect

**Debug**:
```go
inspector := routing.NewMuxInspector(mux)
status, _ := inspector.TestRequest("GET", "/bill")
// Should NOT be 404 if route exists
```

#### 4. Test Timeout

**Problem**: Tests run slowly (> 5s)

**Solution**:
```bash
# Run in parallel
go test -parallel 8 ./cmd/...

# Skip slow tests
go test -short ./cmd/...
```

## Extending for New Services

### Adding a New Service

#### 1. Create WADL File

```xml
<!-- wadl/myservice.wadl -->
<application>
  <resources base="...">
    <resource id="Resource1" wx:samplePath="/resource1">
      <method id="GETResource1" name="GET"/>
      <method id="POSTResource1" name="POST"/>
    </resource>
  </resources>
</application>
```

#### 2. Generate Test Data

```bash
go run ./tools/generate_test_routes.go -service MyService
```

#### 3. Generate Test Template

```bash
./tools/generate_all_service_tests.sh
```

#### 4. Implement Routes

Create `cmd/MyService/main.go` with route registration:

```go
func registerMyServiceRoutes(mux *http.ServeMux, h *Handler) {
    mux.Handle("GET /resource1", http.HandlerFunc(h.GETResource1))
    mux.Handle("POST /resource1", http.HandlerFunc(h.POSTResource1))
    // ... etc
}
```

#### 5. Implement Test

Update `cmd/MyService/myservice_route_registration_test.go`:

```go
func TestMyServiceRouteRegistration(t *testing.T) {
    // Load expectations
    expectations, _ := routing.LoadTestDataByServiceName("MyService")
    
    // Create mux and register routes
    mux := http.NewServeMux()
    h := &Handler{}
    registerMyServiceRoutes(mux, h)
    
    // Verify routes
    helper := routing.NewTestHelper(mux)
    helper.RegisterExpectedRoutes(expectations)
    all, _ := helper.AssertAllRoutesExist(expectations)
    if !all {
        t.Fatal("Routes not registered")
    }
}
```

#### 6. Run Tests

```bash
go test -v ./cmd/MyService/...
```

## Test Maintenance

### When WADL Changes

1. Update WADL file
2. Regenerate test data:
   ```bash
   go run ./tools/generate_test_routes.go
   ```
3. Tests automatically use new expectations

### When Adding Routes

1. Update WADL file with new resources/methods
2. Regenerate test data
3. Update main.go with route registration
4. Tests verify routes are registered

### When Updating Handlers

1. Update handler methods in handler.go
2. Ensure route registration in main.go still works
3. Run tests to verify:
   ```bash
   go test -v ./cmd/<Service>/...
   ```

## Performance

### Expected Test Execution Times

- Per service: < 100ms
- All 23 services: < 3 seconds (parallel)
- With coverage: < 5 seconds

### Optimization Tips

```bash
# Parallel execution (fastest)
go test -parallel 16 ./cmd/...

# Skip coverage (faster)
go test ./cmd/... -count=1

# Run only changed tests
go test ./cmd/Bill/... -run Route
```

## Best Practices

1. **Keep tests focused**: One test per concern
2. **Use fixtures**: Load test data from testdata/
3. **Parallel execution**: Use -parallel flag
4. **Version control**: Commit test data files
5. **CI/CD integration**: Run on every commit
6. **Documentation**: Update this guide with changes

## Summary

| Phase | Purpose | Files | Status |
|-------|---------|-------|--------|
| 1 | Framework design | routing/*.go | ✅ Complete |
| 2 | Route extraction | wadl_extractor.go | ✅ Complete |
| 3 | Route tests | cmd/*_route_registration_test.go | ✅ Complete |
| 4 | HTTP verb tests | cmd/*_http_methods_test.go | 🔄 In Progress |
| 5 | CI/CD integration | Makefile, CI/workflows | ⏳ Ready |
| 6 | Advanced testing | Coverage, performance | ⏳ Ready |

---

## See Also

- `test/routing/README.md` - Framework API reference
- `WADL_RESTORATION_INDEX.md` - Route generation from WADL
- `Makefile` - Test execution targets
