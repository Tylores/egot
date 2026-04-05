// Code generated - DO NOT EDIT
// This test implements route registration tests for Upt service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/UPT/repository/memory"
	"github.com/Tylores/egot/internal/UPT/handler"
	"github.com/Tylores/egot/test/routing"
)

// RegisterUPTRoutes registers all Upt routes on the mux.
func registerUPTRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.Handle("DELETE /upt", http.HandlerFunc(h.DELETEUsagePointList))
	mux.Handle("GET /upt", http.HandlerFunc(h.GETUsagePointList))
	mux.Handle("HEAD /upt", http.HandlerFunc(h.HEADUsagePointList))
	mux.Handle("POST /upt", http.HandlerFunc(h.POSTUsagePointList))
	mux.Handle("PUT /upt", http.HandlerFunc(h.PUTUsagePointList))
	mux.Handle("DELETE /upt/{id1}", http.HandlerFunc(h.DELETEUsagePoint))
	mux.Handle("GET /upt/{id1}", http.HandlerFunc(h.GETUsagePoint))
	mux.Handle("HEAD /upt/{id1}", http.HandlerFunc(h.HEADUsagePoint))
	mux.Handle("POST /upt/{id1}", http.HandlerFunc(h.POSTUsagePoint))
	mux.Handle("PUT /upt/{id1}", http.HandlerFunc(h.PUTUsagePoint))
	mux.Handle("DELETE /upt/{id1}/mr", http.HandlerFunc(h.DELETEMeterReadingList))
	mux.Handle("GET /upt/{id1}/mr", http.HandlerFunc(h.GETMeterReadingList))
	mux.Handle("HEAD /upt/{id1}/mr", http.HandlerFunc(h.HEADMeterReadingList))
	mux.Handle("POST /upt/{id1}/mr", http.HandlerFunc(h.POSTMeterReadingList))
	mux.Handle("PUT /upt/{id1}/mr", http.HandlerFunc(h.PUTMeterReadingList))
	mux.Handle("DELETE /upt/{id1}/mr/{id2}", http.HandlerFunc(h.DELETEMeterReading))
	mux.Handle("GET /upt/{id1}/mr/{id2}", http.HandlerFunc(h.GETMeterReading))
	mux.Handle("HEAD /upt/{id1}/mr/{id2}", http.HandlerFunc(h.HEADMeterReading))
	mux.Handle("POST /upt/{id1}/mr/{id2}", http.HandlerFunc(h.POSTMeterReading))
	mux.Handle("PUT /upt/{id1}/mr/{id2}", http.HandlerFunc(h.PUTMeterReading))
	mux.Handle("DELETE /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.DELETEReadingSetList))
	mux.Handle("GET /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.GETReadingSetList))
	mux.Handle("HEAD /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.HEADReadingSetList))
	mux.Handle("POST /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.POSTReadingSetList))
	mux.Handle("PUT /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.PUTReadingSetList))
	mux.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.DELETEReadingSet))
	mux.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.GETReadingSet))
	mux.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.HEADReadingSet))
	mux.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.POSTReadingSet))
	mux.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.PUTReadingSet))
	mux.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.DELETEReadingList))
	mux.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.GETReadingList))
	mux.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.HEADReadingList))
	mux.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.POSTReadingList))
	mux.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.PUTReadingList))
	mux.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.DELETEReading))
	mux.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.GETReading))
	mux.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.HEADReading))
	mux.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.POSTReading))
	mux.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.PUTReading))
	mux.Handle("DELETE /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.DELETEReadingType))
	mux.Handle("GET /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.GETReadingType))
	mux.Handle("HEAD /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.HEADReadingType))
	mux.Handle("POST /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.POSTReadingType))
	mux.Handle("PUT /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.PUTReadingType))
}

// TestUPTRouteRegistration verifies all Upt routes are correctly registered.
func TestUPTRouteRegistration(t *testing.T) {
	// Create test mux with Upt routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerUPTRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/upt_routes.json")
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

// TestUPTPathParameters verifies path parameter extraction.
func TestUPTPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/upt": {},
		"/upt/{id1}": {"id1"},
		"/upt/{id1}/mr": {"id1"},
		"/upt/{id1}/mr/{id2}": {"id1", "id2"},
		"/upt/{id1}/mr/{id2}/rs": {"id1", "id2"},
		"/upt/{id1}/mr/{id2}/rs/{id3}": {"id1", "id2", "id3"},
		"/upt/{id1}/mr/{id2}/rs/{id3}/r": {"id1", "id2", "id3"},
		"/upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}": {"id1", "id2", "id3", "id4"},
		"/upt/{id1}/mr/{id2}/rt": {"id1", "id2"},
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

// TestUPTHTTPMethods verifies HTTP method support.
func TestUPTHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerUPTRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Load test data for verification
	routeData, err := ioutil.ReadFile("../../test/testdata/upt_routes.json")
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
