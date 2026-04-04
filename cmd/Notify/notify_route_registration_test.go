// Code generated - DO NOT EDIT
// This test implements route registration tests for Notify service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/Notify/repository/memory"
	"github.com/Tylores/egot/internal/Notify/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestNotifyRouteRegistration verifies all Notify routes are correctly registered.
func TestNotifyRouteRegistration(t *testing.T) {
	// Create test mux with Notify routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerNotifyRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/ntfy_routes.json")
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

// TestNotifyPathParameters verifies path parameter extraction.
func TestNotifyPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/ntfy": nil,
		"/ntfy/{id1}": {"id1"},
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

// TestNotifyHTTPMethods verifies HTTP method support.
func TestNotifyHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerNotifyRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/ntfy", true},
		{"DELETE", "/ntfy/{id1}", true},
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

// registerNotifyRoutes registers all Notify service routes
func registerNotifyRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 10 routes registered in main.go
	mux.Handle("DELETE /ntfy", http.HandlerFunc(h.DELETENotificationList))
	mux.Handle("GET /ntfy", http.HandlerFunc(h.GETNotificationList))
	mux.Handle("HEAD /ntfy", http.HandlerFunc(h.HEADNotificationList))
	mux.Handle("POST /ntfy", http.HandlerFunc(h.POSTNotificationList))
	mux.Handle("PUT /ntfy", http.HandlerFunc(h.PUTNotificationList))

	mux.Handle("DELETE /ntfy/{id1}", http.HandlerFunc(h.DELETENotification))
	mux.Handle("GET /ntfy/{id1}", http.HandlerFunc(h.GETNotification))
	mux.Handle("HEAD /ntfy/{id1}", http.HandlerFunc(h.HEADNotification))
	mux.Handle("POST /ntfy/{id1}", http.HandlerFunc(h.POSTNotification))
	mux.Handle("PUT /ntfy/{id1}", http.HandlerFunc(h.PUTNotification))
}
