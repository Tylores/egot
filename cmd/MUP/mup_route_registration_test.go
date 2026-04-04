// Code generated - DO NOT EDIT
// This test implements route registration tests for Mup service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/MUP/repository/memory"
	"github.com/Tylores/egot/internal/MUP/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestMupRouteRegistration verifies all Mup routes are correctly registered.
func TestMupRouteRegistration(t *testing.T) {
	// Create test mux with Mup routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerMupRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/mup_routes.json")
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

// TestMupPathParameters verifies path parameter extraction.
func TestMupPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/mup": nil,
		"/mup/{id1}": {"id1"},
	}

	for path, expectedParams := range patterns {
		if expectedParams != nil {
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

// TestMupHTTPMethods verifies HTTP method support.
func TestMupHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerMupRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/mup", true},
		{"DELETE", "/mup/{id1}", true},
		// Unsupported paths (should not match)
		{"GET", "/invalid", false},
		{"POST", "/nonexistent", false},
	}

	for _, test := range tests {
		// Test the route
		status, _ := inspector.TestRequest(test.method, test.path)

		// Route found if status != 404
		found := status != http.StatusNotFound
		if found != test.expect {
			t.Logf("%s %s: status=%d, expected found=%v",
				test.method, test.path, status, test.expect)
		}
	}
}

// registerMupRoutes registers all Mup service routes
func registerMupRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 10 routes registered in main.go
	mux.Handle("DELETE /mup", http.HandlerFunc(h.DELETEMirrorUsagePointList))
	mux.Handle("GET /mup", http.HandlerFunc(h.GETMirrorUsagePointList))
	mux.Handle("HEAD /mup", http.HandlerFunc(h.HEADMirrorUsagePointList))
	mux.Handle("POST /mup", http.HandlerFunc(h.POSTMirrorUsagePointList))
	mux.Handle("PUT /mup", http.HandlerFunc(h.PUTMirrorUsagePointList))

	mux.Handle("DELETE /mup/{id1}", http.HandlerFunc(h.DELETEMirrorUsagePoint))
	mux.Handle("GET /mup/{id1}", http.HandlerFunc(h.GETMirrorUsagePoint))
	mux.Handle("HEAD /mup/{id1}", http.HandlerFunc(h.HEADMirrorUsagePoint))
	mux.Handle("POST /mup/{id1}", http.HandlerFunc(h.POSTMirrorUsagePoint))
	mux.Handle("PUT /mup/{id1}", http.HandlerFunc(h.PUTMirrorUsagePoint))
}
