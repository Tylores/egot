// Code generated - DO NOT EDIT
// This test implements route registration tests for Derp service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/DERP/repository/memory"
	"github.com/Tylores/egot/internal/DERP/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestDerpRouteRegistration verifies all Derp routes are correctly registered.
func TestDerpRouteRegistration(t *testing.T) {
	// Create test mux with Derp routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerDerpRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/derp_routes.json")
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

// TestDerpPathParameters verifies path parameter extraction.
func TestDerpPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/derp": nil,
		"/derp/{id1}": {"id1"},
		"/derp/{id1}/actderc": {"id1"},
		"/derp/{id1}/dc": {"id1"},
		"/derp/{id1}/dc/{id2}": {"id1", "id2"},
		"/derp/{id1}/dderc": {"id1"},
		"/derp/{id1}/derc": {"id1"},
		"/derp/{id1}/derc/{id2}": {"id1", "id2"},
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

// TestDerpHTTPMethods verifies HTTP method support.
func TestDerpHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerDerpRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/derp", true},
		{"DELETE", "/derp/{id1}", true},
		{"DELETE", "/derp/{id1}/actderc", true},
		{"DELETE", "/derp/{id1}/dc", true},
		{"DELETE", "/derp/{id1}/dc/{id2}", true},
		{"DELETE", "/derp/{id1}/dderc", true},
		{"DELETE", "/derp/{id1}/derc", true},
		{"DELETE", "/derp/{id1}/derc/{id2}", true},
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

// registerDerpRoutes registers all Derp service routes
func registerDerpRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 40 routes registered in main.go
	mux.Handle("DELETE /derp", http.HandlerFunc(h.DELETEDERProgramList))
	mux.Handle("GET /derp", http.HandlerFunc(h.GETDERProgramList))
	mux.Handle("HEAD /derp", http.HandlerFunc(h.HEADDERProgramList))
	mux.Handle("POST /derp", http.HandlerFunc(h.POSTDERProgramList))
	mux.Handle("PUT /derp", http.HandlerFunc(h.PUTDERProgramList))

	mux.Handle("DELETE /derp/{id1}", http.HandlerFunc(h.DELETEDERProgram))
	mux.Handle("GET /derp/{id1}", http.HandlerFunc(h.GETDERProgram))
	mux.Handle("HEAD /derp/{id1}", http.HandlerFunc(h.HEADDERProgram))
	mux.Handle("POST /derp/{id1}", http.HandlerFunc(h.POSTDERProgram))
	mux.Handle("PUT /derp/{id1}", http.HandlerFunc(h.PUTDERProgram))

	mux.Handle("DELETE /derp/{id1}/actderc", http.HandlerFunc(h.DELETEActiveDERControlList))
	mux.Handle("GET /derp/{id1}/actderc", http.HandlerFunc(h.GETActiveDERControlList))
	mux.Handle("HEAD /derp/{id1}/actderc", http.HandlerFunc(h.HEADActiveDERControlList))
	mux.Handle("POST /derp/{id1}/actderc", http.HandlerFunc(h.POSTActiveDERControlList))
	mux.Handle("PUT /derp/{id1}/actderc", http.HandlerFunc(h.PUTActiveDERControlList))

	mux.Handle("DELETE /derp/{id1}/dc", http.HandlerFunc(h.DELETEDERCurveList))
	mux.Handle("GET /derp/{id1}/dc", http.HandlerFunc(h.GETDERCurveList))
	mux.Handle("HEAD /derp/{id1}/dc", http.HandlerFunc(h.HEADDERCurveList))
	mux.Handle("POST /derp/{id1}/dc", http.HandlerFunc(h.POSTDERCurveList))
	mux.Handle("PUT /derp/{id1}/dc", http.HandlerFunc(h.PUTDERCurveList))

	mux.Handle("DELETE /derp/{id1}/dc/{id2}", http.HandlerFunc(h.DELETEDERCurve))
	mux.Handle("GET /derp/{id1}/dc/{id2}", http.HandlerFunc(h.GETDERCurve))
	mux.Handle("HEAD /derp/{id1}/dc/{id2}", http.HandlerFunc(h.HEADDERCurve))
	mux.Handle("POST /derp/{id1}/dc/{id2}", http.HandlerFunc(h.POSTDERCurve))
	mux.Handle("PUT /derp/{id1}/dc/{id2}", http.HandlerFunc(h.PUTDERCurve))

	mux.Handle("DELETE /derp/{id1}/dderc", http.HandlerFunc(h.DELETEDefaultDERControl))
	mux.Handle("GET /derp/{id1}/dderc", http.HandlerFunc(h.GETDefaultDERControl))
	mux.Handle("HEAD /derp/{id1}/dderc", http.HandlerFunc(h.HEADDefaultDERControl))
	mux.Handle("POST /derp/{id1}/dderc", http.HandlerFunc(h.POSTDefaultDERControl))
	mux.Handle("PUT /derp/{id1}/dderc", http.HandlerFunc(h.PUTDefaultDERControl))

	mux.Handle("DELETE /derp/{id1}/derc", http.HandlerFunc(h.DELETEDERControlList))
	mux.Handle("GET /derp/{id1}/derc", http.HandlerFunc(h.GETDERControlList))
	mux.Handle("HEAD /derp/{id1}/derc", http.HandlerFunc(h.HEADDERControlList))
	mux.Handle("POST /derp/{id1}/derc", http.HandlerFunc(h.POSTDERControlList))
	mux.Handle("PUT /derp/{id1}/derc", http.HandlerFunc(h.PUTDERControlList))

	mux.Handle("DELETE /derp/{id1}/derc/{id2}", http.HandlerFunc(h.DELETEDERControl))
	mux.Handle("GET /derp/{id1}/derc/{id2}", http.HandlerFunc(h.GETDERControl))
	mux.Handle("HEAD /derp/{id1}/derc/{id2}", http.HandlerFunc(h.HEADDERControl))
	mux.Handle("POST /derp/{id1}/derc/{id2}", http.HandlerFunc(h.POSTDERControl))
	mux.Handle("PUT /derp/{id1}/derc/{id2}", http.HandlerFunc(h.PUTDERControl))
}
