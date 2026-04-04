# Test Routing Framework

This package provides utilities for testing HTTP route registration in microservices.

## Overview

The `routing` package includes tools to:
- Validate route registration (method + path)
- Verify HTTP method support
- Extract and validate path parameters
- Compare expected vs. actual routes
- Generate route expectations from WADL files

## Components

### RouteValidator
Core validator for route registration verification.

```go
validator := routing.NewRouteValidator(mux)
validator.RegisterRoute("GET", "/bill")
validator.RegisterRoute("POST", "/bill/{id1}/ca/{id2}")

// Assert routes exist
all := validator.AssertAllRoutesExist(expectations)
count := validator.GetRegisteredCount()
```

### HTTPMethodValidator
Validates HTTP method support for routes.

```go
methods := routing.NewHTTPMethodValidator()
methods.AddSupportedMethod("/bill", "GET")
methods.AddSupportedMethod("/bill", "POST")

// Check if method is supported
supported := methods.IsSupportedMethod("/bill", "GET")
```

### PathParameterValidator
Extracts and validates path parameters.

```go
params := routing.NewPathParameterValidator()
params.AddPathPattern("/bill/{id1}/ca/{id2}", []string{"id1", "id2"})

// Extract parameters from path
extracted := params.ExtractParameterNames("/bill/{id1}/ca/{id2}")
// Returns: ["id1", "id2"]
```

### MuxInspector
Inspects HTTP mux for registered routes.

```go
inspector := routing.NewMuxInspector(mux)

// Discover routes by testing patterns
discovered := inspector.DiscoverRoutes(testPaths)

// Test individual routes
status, _ := inspector.TestRequest("GET", "/bill")
matched := inspector.AssertRouteMatchesMethod("GET", "/bill")
```

### TestHelper
Convenience wrapper combining all validators.

```go
helper := routing.NewTestHelper(mux)
helper.RegisterExpectedRoute("GET", "/bill")
helper.RegisterExpectedRoute("POST", "/bill/{id1}/ca/{id2}")

// Get summary
summary := helper.GetSummary(expectations)
```

### RouteTestBuilder
Fluent API for constructing test scenarios.

```go
builder := routing.NewRouteTestBuilder()
builder.
    AddRoutes([]string{"GET", "POST"}, []string{"/bill", "/ca"}).
    WithParameters("/bill/{id1}", []string{"id1"})

expectations, methods, params := builder.Build()
```

## Usage Patterns

### Basic Route Testing

```go
func TestBillRoutes(t *testing.T) {
    // Create expected routes
    expectations := []routing.RouteExpectation{
        {Method: "GET", Path: "/bill"},
        {Method: "POST", Path: "/bill/{id1}"},
        {Method: "DELETE", Path: "/bill/{id1}"},
    }
    
    // Create test helper
    helper := routing.NewTestHelper(mux)
    helper.RegisterExpectedRoutes(expectations)
    
    // Verify all routes exist
    all, missing := helper.AssertAllRoutesExist(expectations)
    if !all {
        t.Errorf("Missing routes: %v", missing)
    }
}
```

### Path Parameter Testing

```go
func TestPathParameters(t *testing.T) {
    helper := routing.NewTestHelper(mux)
    
    path := "/bill/{id1}/ca/{id2}"
    expected := []string{"id1", "id2"}
    
    valid := helper.VerifyPathParameters(path, expected)
    if !valid {
        t.Errorf("Path parameters invalid for %s", path)
    }
}
```

### HTTP Method Testing

```go
func TestHTTPMethods(t *testing.T) {
    inspector := routing.NewMuxInspector(mux)
    
    tests := []struct {
        method string
        path   string
        expect bool
    }{
        {"GET", "/bill", true},
        {"POST", "/bill", true},
        {"PATCH", "/bill", false}, // Not supported
    }
    
    for _, test := range tests {
        matched := inspector.AssertRouteMatchesMethod(test.method, test.path)
        if matched != test.expect {
            t.Errorf("%s %s: expected %v, got %v", 
                test.method, test.path, test.expect, matched)
        }
    }
}
```

## Route Expectation Format

Routes are specified as `RouteExpectation` structs:

```go
type RouteExpectation struct {
    Method string // HTTP method: GET, POST, PUT, DELETE, HEAD
    Path   string // URL path: /bill, /bill/{id1}, /bill/{id1}/ca/{id2}
}
```

## Path Parameter Syntax

Path parameters use standard Go mux syntax with braces:

- Single parameter: `/bill/{id1}`
- Multiple parameters: `/bill/{id1}/ca/{id2}/bp/{id3}`
- Parameter names are typically numeric: `{id1}`, `{id2}`, `{id3}`

## Common Test Patterns

### Verify All Routes Count

```go
count, registered, expected := helper.AssertRouteCount(expectations)
if !count {
    t.Errorf("Route count mismatch: got %d, expected %d", registered, expected)
}
```

### Get Test Summary

```go
summary := helper.GetSummary(expectations)
fmt.Printf("Expected: %d, Registered: %d, Missing: %d\n",
    summary.TotalExpected,
    summary.TotalRegistered,
    summary.MissingCount)
```

### Test with Builder

```go
builder := routing.NewRouteTestBuilder()
builder.AddPath("/bill").     // Adds all 4 HTTP verbs
        AddPath("/ca").
        WithParameters("/bill/{id1}", []string{"id1"})

expectations := builder.GetExpectations()
```

## Integration with Service Tests

Each microservice test file uses this framework:

1. Create RouteExpectation list from WADL expectations
2. Create TestHelper with service mux
3. Register expected routes
4. Assert all routes exist and counts match
5. Test HTTP methods and path parameters

Example: `cmd/Bill/bill_route_registration_test.go`

```go
package main

import (
    "testing"
    "github.com/Tylores/egot/test/routing"
)

func TestBillRouteRegistration(t *testing.T) {
    mux := http.NewServeMux()
    registerBillRoutes(mux, &BillHandler{})
    
    expectations := []routing.RouteExpectation{
        {Method: "GET", Path: "/bill"},
        // ... 69 more routes from WADL
    }
    
    helper := routing.NewTestHelper(mux)
    helper.RegisterExpectedRoutes(expectations)
    
    all, missing := helper.AssertAllRoutesExist(expectations)
    if !all {
        t.Fatalf("Missing %d routes", len(missing))
    }
}
```

## Troubleshooting

### Routes Not Found

1. Verify route path matches exactly (case-sensitive)
2. Check HTTP method is uppercase: GET, POST, PUT, DELETE
3. Ensure path parameters use correct syntax: `{paramName}`

### Test Failures

1. Get route summary: `fmt.Printf("%v", helper.GetSummary(expectations))`
2. Check for duplicate routes: Use `inspector.GetRegisteredRoutes()`
3. Compare expected vs. actual: `result := inspector.CompareRoutes(...)`

### Parameter Issues

1. Extract parameters: `params := helper.ExtractPathParameters(path)`
2. Verify pattern: `valid := helper.AssertPathHasParameters(path)`
3. Check parameter names match WADL: `expected := []string{"id1", "id2", "id3"}`

## Testing Checklist

- [ ] All routes from WADL are registered
- [ ] Route count matches WADL resource count
- [ ] HTTP methods match WADL definitions
- [ ] Path parameters are correctly extracted
- [ ] No duplicate routes
- [ ] No unexpected extra routes

## See Also

- `WADL_RESTORATION_INDEX.md` - Route generation from WADL
- `TEST_ROUTING_GUIDE.md` - Comprehensive testing guide
- Service test files: `cmd/<Service>/<Service>_route_registration_test.go`
