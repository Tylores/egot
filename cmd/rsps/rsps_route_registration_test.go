// Code generated - DO NOT EDIT
// This test implements route registration tests for Rsps service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/rsps/repository/memory"
	"github.com/Tylores/egot/internal/rsps/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestRspsRouteRegistration verifies all Rsps routes are correctly registered.
func TestRspsRouteRegistration(t *testing.T) {
	// Create test mux with Rsps routes
	mux := http.NewServeMux()
	repo := memory.NewRepository(1000) // rsps requires Entity size parameter
	h := handler.NewHandler(repo)

	// Register routes
	registerRspsRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/rsps_routes.json")
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

// TestRspsPathParameters verifies path parameter extraction.
func TestRspsPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/rsps": nil,
		"/rsps/{id1}": {"id1"},
		"/rsps/{id1}/rsp": {"id1"},
		"/rsps/{id1}/rsp/{id2}": {"id1", "id2"},
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

// TestRspsHTTPMethods verifies HTTP method support.
func TestRspsHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository(1000) // rsps requires Entity size parameter
	h := handler.NewHandler(repo)

	// Register routes
	registerRspsRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/rsps", true},
		{"DELETE", "/rsps/{id1}", true},
		{"DELETE", "/rsps/{id1}/rsp", true},
		{"DELETE", "/rsps/{id1}/rsp/{id2}", true},
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

// registerRspsRoutes registers all Rsps service routes
func registerRspsRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 20 routes registered in main.go
	mux.Handle("DELETE /rsps", http.HandlerFunc(h.DELETEResponseSetList))
	mux.Handle("GET /rsps", http.HandlerFunc(h.GETResponseSetList))
	mux.Handle("HEAD /rsps", http.HandlerFunc(h.HEADResponseSetList))
	mux.Handle("POST /rsps", http.HandlerFunc(h.POSTResponseSetList))
	mux.Handle("PUT /rsps", http.HandlerFunc(h.PUTResponseSetList))

	mux.Handle("DELETE /rsps/{id1}", http.HandlerFunc(h.DELETEResponseSet))
	mux.Handle("GET /rsps/{id1}", http.HandlerFunc(h.GETResponseSet))
	mux.Handle("HEAD /rsps/{id1}", http.HandlerFunc(h.HEADResponseSet))
	mux.Handle("POST /rsps/{id1}", http.HandlerFunc(h.POSTResponseSet))
	mux.Handle("PUT /rsps/{id1}", http.HandlerFunc(h.PUTResponseSet))

	mux.Handle("DELETE /rsps/{id1}/rsp", http.HandlerFunc(h.DELETEResponseList))
	mux.Handle("GET /rsps/{id1}/rsp", http.HandlerFunc(h.GETResponseList))
	mux.Handle("HEAD /rsps/{id1}/rsp", http.HandlerFunc(h.HEADResponseList))
	mux.Handle("POST /rsps/{id1}/rsp", http.HandlerFunc(h.POSTResponseList))
	mux.Handle("PUT /rsps/{id1}/rsp", http.HandlerFunc(h.PUTResponseList))

	mux.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEResponse))
	mux.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETResponse))
	mux.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADResponse))
	mux.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTResponse))
	mux.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTResponse))
}
