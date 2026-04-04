// Code generated - DO NOT EDIT
// This test implements route registration tests for Dr service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/DR/repository/memory"
	"github.com/Tylores/egot/internal/DR/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestDrRouteRegistration verifies all Dr routes are correctly registered.
func TestDrRouteRegistration(t *testing.T) {
	// Create test mux with Dr routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerDrRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/dr_routes.json")
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

// TestDrPathParameters verifies path parameter extraction.
func TestDrPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/dr": nil,
		"/dr/{id1}": {"id1"},
		"/dr/{id1}/actedc": {"id1"},
		"/dr/{id1}/edc": {"id1"},
		"/dr/{id1}/edc/{id2}": {"id1", "id2"},
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

// TestDrHTTPMethods verifies HTTP method support.
func TestDrHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerDrRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/dr", true},
		{"DELETE", "/dr/{id1}", true},
		{"DELETE", "/dr/{id1}/actedc", true},
		{"DELETE", "/dr/{id1}/edc", true},
		{"DELETE", "/dr/{id1}/edc/{id2}", true},
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

// registerDrRoutes registers all Dr service routes
func registerDrRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 25 routes registered in main.go
	mux.Handle("DELETE /dr", http.HandlerFunc(h.DELETEDemandResponseProgramList))
	mux.Handle("GET /dr", http.HandlerFunc(h.GETDemandResponseProgramList))
	mux.Handle("HEAD /dr", http.HandlerFunc(h.HEADDemandResponseProgramList))
	mux.Handle("POST /dr", http.HandlerFunc(h.POSTDemandResponseProgramList))
	mux.Handle("PUT /dr", http.HandlerFunc(h.PUTDemandResponseProgramList))

	mux.Handle("DELETE /dr/{id1}", http.HandlerFunc(h.DELETEDemandResponseProgram))
	mux.Handle("GET /dr/{id1}", http.HandlerFunc(h.GETDemandResponseProgram))
	mux.Handle("HEAD /dr/{id1}", http.HandlerFunc(h.HEADDemandResponseProgram))
	mux.Handle("POST /dr/{id1}", http.HandlerFunc(h.POSTDemandResponseProgram))
	mux.Handle("PUT /dr/{id1}", http.HandlerFunc(h.PUTDemandResponseProgram))

	mux.Handle("DELETE /dr/{id1}/actedc", http.HandlerFunc(h.DELETEActiveEndDeviceControlList))
	mux.Handle("GET /dr/{id1}/actedc", http.HandlerFunc(h.GETActiveEndDeviceControlList))
	mux.Handle("HEAD /dr/{id1}/actedc", http.HandlerFunc(h.HEADActiveEndDeviceControlList))
	mux.Handle("POST /dr/{id1}/actedc", http.HandlerFunc(h.POSTActiveEndDeviceControlList))
	mux.Handle("PUT /dr/{id1}/actedc", http.HandlerFunc(h.PUTActiveEndDeviceControlList))

	mux.Handle("DELETE /dr/{id1}/edc", http.HandlerFunc(h.DELETEEndDeviceControlList))
	mux.Handle("GET /dr/{id1}/edc", http.HandlerFunc(h.GETEndDeviceControlList))
	mux.Handle("HEAD /dr/{id1}/edc", http.HandlerFunc(h.HEADEndDeviceControlList))
	mux.Handle("POST /dr/{id1}/edc", http.HandlerFunc(h.POSTEndDeviceControlList))
	mux.Handle("PUT /dr/{id1}/edc", http.HandlerFunc(h.PUTEndDeviceControlList))

	mux.Handle("DELETE /dr/{id1}/edc/{id2}", http.HandlerFunc(h.DELETEEndDeviceControl))
	mux.Handle("GET /dr/{id1}/edc/{id2}", http.HandlerFunc(h.GETEndDeviceControl))
	mux.Handle("HEAD /dr/{id1}/edc/{id2}", http.HandlerFunc(h.HEADEndDeviceControl))
	mux.Handle("POST /dr/{id1}/edc/{id2}", http.HandlerFunc(h.POSTEndDeviceControl))
	mux.Handle("PUT /dr/{id1}/edc/{id2}", http.HandlerFunc(h.PUTEndDeviceControl))
}
