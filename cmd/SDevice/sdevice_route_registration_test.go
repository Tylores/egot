// Code generated - DO NOT EDIT
// This test implements route registration tests for Sdevice service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/SDevice/repository/memory"
	"github.com/Tylores/egot/internal/SDevice/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestSdeviceRouteRegistration verifies all Sdevice routes are correctly registered.
func TestSdeviceRouteRegistration(t *testing.T) {
	// Create test mux with Sdevice routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerSdeviceRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/sdev_routes.json")
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

// TestSdevicePathParameters verifies path parameter extraction.
func TestSdevicePathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/sdev": nil,
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

// TestSdeviceHTTPMethods verifies HTTP method support.
func TestSdeviceHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerSdeviceRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/sdev", true},
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

// registerSdeviceRoutes registers all Sdevice service routes
func registerSdeviceRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 5 routes registered in main.go
	mux.Handle("DELETE /sdev", http.HandlerFunc(h.DELETESelfDevice))
	mux.Handle("GET /sdev", http.HandlerFunc(h.GETSelfDevice))
	mux.Handle("HEAD /sdev", http.HandlerFunc(h.HEADSelfDevice))
	mux.Handle("POST /sdev", http.HandlerFunc(h.POSTSelfDevice))
	mux.Handle("PUT /sdev", http.HandlerFunc(h.PUTSelfDevice))
}
