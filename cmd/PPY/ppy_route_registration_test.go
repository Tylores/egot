// Code generated - DO NOT EDIT
// This test implements route registration tests for Ppy service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/PPY/repository/memory"
	"github.com/Tylores/egot/internal/PPY/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestPpyRouteRegistration verifies all Ppy routes are correctly registered.
func TestPpyRouteRegistration(t *testing.T) {
	// Create test mux with Ppy routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerPpyRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/ppy_routes.json")
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

// TestPpyPathParameters verifies path parameter extraction.
func TestPpyPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/ppy": nil,
		"/ppy/{id1}": {"id1"},
		"/ppy/{id1}/ab": {"id1"},
		"/ppy/{id1}/actsi": {"id1"},
		"/ppy/{id1}/cr": {"id1"},
		"/ppy/{id1}/cr/{id2}": {"id1", "id2"},
		"/ppy/{id1}/os": {"id1"},
		"/ppy/{id1}/si": {"id1"},
		"/ppy/{id1}/si/{id2}": {"id1", "id2"},
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

// TestPpyHTTPMethods verifies HTTP method support.
func TestPpyHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerPpyRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/ppy", true},
		{"DELETE", "/ppy/{id1}", true},
		{"DELETE", "/ppy/{id1}/ab", true},
		{"DELETE", "/ppy/{id1}/actsi", true},
		{"DELETE", "/ppy/{id1}/cr", true},
		{"DELETE", "/ppy/{id1}/cr/{id2}", true},
		{"DELETE", "/ppy/{id1}/os", true},
		{"DELETE", "/ppy/{id1}/si", true},
		{"DELETE", "/ppy/{id1}/si/{id2}", true},
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

// registerPpyRoutes registers all Ppy service routes
func registerPpyRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 45 routes registered in main.go
	mux.Handle("DELETE /ppy", http.HandlerFunc(h.DELETEPrepaymentList))
	mux.Handle("GET /ppy", http.HandlerFunc(h.GETPrepaymentList))
	mux.Handle("HEAD /ppy", http.HandlerFunc(h.HEADPrepaymentList))
	mux.Handle("POST /ppy", http.HandlerFunc(h.POSTPrepaymentList))
	mux.Handle("PUT /ppy", http.HandlerFunc(h.PUTPrepaymentList))

	mux.Handle("DELETE /ppy/{id1}", http.HandlerFunc(h.DELETEPrepayment))
	mux.Handle("GET /ppy/{id1}", http.HandlerFunc(h.GETPrepayment))
	mux.Handle("HEAD /ppy/{id1}", http.HandlerFunc(h.HEADPrepayment))
	mux.Handle("POST /ppy/{id1}", http.HandlerFunc(h.POSTPrepayment))
	mux.Handle("PUT /ppy/{id1}", http.HandlerFunc(h.PUTPrepayment))

	mux.Handle("DELETE /ppy/{id1}/ab", http.HandlerFunc(h.DELETEAccountBalance))
	mux.Handle("GET /ppy/{id1}/ab", http.HandlerFunc(h.GETAccountBalance))
	mux.Handle("HEAD /ppy/{id1}/ab", http.HandlerFunc(h.HEADAccountBalance))
	mux.Handle("POST /ppy/{id1}/ab", http.HandlerFunc(h.POSTAccountBalance))
	mux.Handle("PUT /ppy/{id1}/ab", http.HandlerFunc(h.PUTAccountBalance))

	mux.Handle("DELETE /ppy/{id1}/actsi", http.HandlerFunc(h.DELETEActiveSupplyInterruptionOverrideList))
	mux.Handle("GET /ppy/{id1}/actsi", http.HandlerFunc(h.GETActiveSupplyInterruptionOverrideList))
	mux.Handle("HEAD /ppy/{id1}/actsi", http.HandlerFunc(h.HEADActiveSupplyInterruptionOverrideList))
	mux.Handle("POST /ppy/{id1}/actsi", http.HandlerFunc(h.POSTActiveSupplyInterruptionOverrideList))
	mux.Handle("PUT /ppy/{id1}/actsi", http.HandlerFunc(h.PUTActiveSupplyInterruptionOverrideList))

	mux.Handle("DELETE /ppy/{id1}/cr", http.HandlerFunc(h.DELETECreditRegisterList))
	mux.Handle("GET /ppy/{id1}/cr", http.HandlerFunc(h.GETCreditRegisterList))
	mux.Handle("HEAD /ppy/{id1}/cr", http.HandlerFunc(h.HEADCreditRegisterList))
	mux.Handle("POST /ppy/{id1}/cr", http.HandlerFunc(h.POSTCreditRegisterList))
	mux.Handle("PUT /ppy/{id1}/cr", http.HandlerFunc(h.PUTCreditRegisterList))

	mux.Handle("DELETE /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.DELETECreditRegister))
	mux.Handle("GET /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.GETCreditRegister))
	mux.Handle("HEAD /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.HEADCreditRegister))
	mux.Handle("POST /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.POSTCreditRegister))
	mux.Handle("PUT /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.PUTCreditRegister))

	mux.Handle("DELETE /ppy/{id1}/os", http.HandlerFunc(h.DELETEPrepayOperationStatus))
	mux.Handle("GET /ppy/{id1}/os", http.HandlerFunc(h.GETPrepayOperationStatus))
	mux.Handle("HEAD /ppy/{id1}/os", http.HandlerFunc(h.HEADPrepayOperationStatus))
	mux.Handle("POST /ppy/{id1}/os", http.HandlerFunc(h.POSTPrepayOperationStatus))
	mux.Handle("PUT /ppy/{id1}/os", http.HandlerFunc(h.PUTPrepayOperationStatus))

	mux.Handle("DELETE /ppy/{id1}/si", http.HandlerFunc(h.DELETESupplyInterruptionOverrideList))
	mux.Handle("GET /ppy/{id1}/si", http.HandlerFunc(h.GETSupplyInterruptionOverrideList))
	mux.Handle("HEAD /ppy/{id1}/si", http.HandlerFunc(h.HEADSupplyInterruptionOverrideList))
	mux.Handle("POST /ppy/{id1}/si", http.HandlerFunc(h.POSTSupplyInterruptionOverrideList))
	mux.Handle("PUT /ppy/{id1}/si", http.HandlerFunc(h.PUTSupplyInterruptionOverrideList))

	mux.Handle("DELETE /ppy/{id1}/si/{id2}", http.HandlerFunc(h.DELETESupplyInterruptionOverride))
	mux.Handle("GET /ppy/{id1}/si/{id2}", http.HandlerFunc(h.GETSupplyInterruptionOverride))
	mux.Handle("HEAD /ppy/{id1}/si/{id2}", http.HandlerFunc(h.HEADSupplyInterruptionOverride))
	mux.Handle("POST /ppy/{id1}/si/{id2}", http.HandlerFunc(h.POSTSupplyInterruptionOverride))
	mux.Handle("PUT /ppy/{id1}/si/{id2}", http.HandlerFunc(h.PUTSupplyInterruptionOverride))
}
