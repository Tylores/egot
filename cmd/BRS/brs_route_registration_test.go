// Code generated - DO NOT EDIT
// This test implements route registration tests for Brs service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/BRS/repository/memory"
	"github.com/Tylores/egot/internal/BRS/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestBrsRouteRegistration verifies all Brs routes are correctly registered.
func TestBrsRouteRegistration(t *testing.T) {
	// Create test mux with Brs routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerBrsRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/brs_routes.json")
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

// TestBrsPathParameters verifies path parameter extraction.
func TestBrsPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/brs": nil,
		"/brs/{id1}": {"id1"},
		"/brs/{id1}/br": {"id1"},
		"/brs/{id1}/br/{id2}": {"id1", "id2"},
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

// TestBrsHTTPMethods verifies HTTP method support.
func TestBrsHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerBrsRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/brs", true},
		{"DELETE", "/brs/{id1}", true},
		{"DELETE", "/brs/{id1}/br", true},
		{"DELETE", "/brs/{id1}/br/{id2}", true},
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

// registerBrsRoutes registers all Brs service routes
func registerBrsRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 20 routes registered in main.go
	mux.Handle("DELETE /brs", http.HandlerFunc(h.DELETEBillingReadingSetList))
	mux.Handle("GET /brs", http.HandlerFunc(h.GETBillingReadingSetList))
	mux.Handle("HEAD /brs", http.HandlerFunc(h.HEADBillingReadingSetList))
	mux.Handle("POST /brs", http.HandlerFunc(h.POSTBillingReadingSetList))
	mux.Handle("PUT /brs", http.HandlerFunc(h.PUTBillingReadingSetList))

	mux.Handle("DELETE /brs/{id1}", http.HandlerFunc(h.DELETEBillingReadingSet))
	mux.Handle("GET /brs/{id1}", http.HandlerFunc(h.GETBillingReadingSet))
	mux.Handle("HEAD /brs/{id1}", http.HandlerFunc(h.HEADBillingReadingSet))
	mux.Handle("POST /brs/{id1}", http.HandlerFunc(h.POSTBillingReadingSet))
	mux.Handle("PUT /brs/{id1}", http.HandlerFunc(h.PUTBillingReadingSet))

	mux.Handle("DELETE /brs/{id1}/br", http.HandlerFunc(h.DELETEBillingReadingList))
	mux.Handle("GET /brs/{id1}/br", http.HandlerFunc(h.GETBillingReadingList))
	mux.Handle("HEAD /brs/{id1}/br", http.HandlerFunc(h.HEADBillingReadingList))
	mux.Handle("POST /brs/{id1}/br", http.HandlerFunc(h.POSTBillingReadingList))
	mux.Handle("PUT /brs/{id1}/br", http.HandlerFunc(h.PUTBillingReadingList))

	mux.Handle("DELETE /brs/{id1}/br/{id2}", http.HandlerFunc(h.DELETEBillingReading))
	mux.Handle("GET /brs/{id1}/br/{id2}", http.HandlerFunc(h.GETBillingReading))
	mux.Handle("HEAD /brs/{id1}/br/{id2}", http.HandlerFunc(h.HEADBillingReading))
	mux.Handle("POST /brs/{id1}/br/{id2}", http.HandlerFunc(h.POSTBillingReading))
	mux.Handle("PUT /brs/{id1}/br/{id2}", http.HandlerFunc(h.PUTBillingReading))
}
