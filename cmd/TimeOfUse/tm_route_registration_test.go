// Code generated - DO NOT EDIT
// This test implements route registration tests for Tm service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/TimeOfUse/repository/memory"
	"github.com/Tylores/egot/internal/TimeOfUse/handler"
	"github.com/Tylores/egot/test/routing"
)

// RegisterTMRoutes registers all Tm routes on the mux.
func registerTMRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.Handle("DELETE /tm", http.HandlerFunc(h.DELETETime))
	mux.Handle("GET /tm", http.HandlerFunc(h.GETTime))
	mux.Handle("HEAD /tm", http.HandlerFunc(h.HEADTime))
	mux.Handle("POST /tm", http.HandlerFunc(h.POSTTime))
	mux.Handle("PUT /tm", http.HandlerFunc(h.PUTTime))
}

// TestTMRouteRegistration verifies all Tm routes are correctly registered.
func TestTMRouteRegistration(t *testing.T) {
	// Create test mux with Tm routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerTMRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/tm_routes.json")
	if err != nil {
		t.Logf("Note: Could not load test data: %v", err)
		t.Skip("Test data not available")
	}

	var testData struct {
		Routes []struct {
			Method string
			Path   string
		}
	}

	if err := json.Unmarshal(routeData, &testData); err != nil {
		t.Fatalf("Failed to parse test data: %v", err)
	}

	// Convert to RouteExpectation format
	expectations := make([]routing.RouteExpectation, 0, len(testData.Routes))
	for _, r := range testData.Routes {
		expectations = append(expectations, routing.RouteExpectation{
			Method: r.Method,
			Path:   r.Path,
		})
	}

	// Create test helper
	helper := routing.NewTestHelper(mux)
	helper.RegisterExpectedRoutes(expectations)

	// Verify routes exist
	all, missing := helper.AssertAllRoutesExist(expectations)
	if !all {
		t.Errorf("Missing %d routes:", len(missing))
		for _, route := range missing {
			t.Errorf("  - %s", route)
		}
	}

	// Verify route count
	matches, registered, expected := helper.AssertRouteCount(expectations)
	if !matches {
		t.Errorf("Route count mismatch: expected %d, got %d", expected, registered)
	}
}

// TestTMPathParameters verifies path parameter extraction.
func TestTMPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/tm": {},
}

	for path, expectedParams := range patterns {
		if len(expectedParams) > 0 {
			validator.AddPathPattern(path, expectedParams)
		}

		// Extract parameters
		extracted := validator.ExtractParameterNames(path)

		// Verify count
		if len(extracted) != len(expectedParams) {
			t.Errorf("Path %s: expected %d params, got %d",
				path, len(expectedParams), len(extracted))
			continue
		}

		// Verify names
		for i, expected := range expectedParams {
			if i < len(extracted) && extracted[i] != expected {
				t.Errorf("Path %s param %d: expected '%s', got '%s'",
					path, i, expected, extracted[i])
			}
		}
	}
}

// TestTMHTTPMethods verifies HTTP method support.
func TestTMHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerTMRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Load test data for verification
	routeData, err := ioutil.ReadFile("../../test/testdata/tm_routes.json")
	if err != nil {
		t.Skip("Test data not available")
	}

	var testData struct {
		Routes []struct {
			Method string
			Path   string
		}
	}

	if err := json.Unmarshal(routeData, &testData); err != nil {
		t.Skip("Invalid test data")
	}

	// Test routes from test data
	for _, route := range testData.Routes {
		status, _ := inspector.TestRequest(route.Method, route.Path)
		if status == http.StatusNotFound {
			t.Errorf("Route not properly registered: %s %s", route.Method, route.Path)
		}
	}
}
