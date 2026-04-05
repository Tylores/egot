// Code generated - DO NOT EDIT
// This test implements route registration tests for Tp service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/TariffProfile/repository/memory"
	"github.com/Tylores/egot/internal/TariffProfile/handler"
	"github.com/Tylores/egot/test/routing"
)

// RegisterTPRoutes registers all Tp routes on the mux.
func registerTPRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.Handle("DELETE /tp", http.HandlerFunc(h.DELETETariffProfileList))
	mux.Handle("GET /tp", http.HandlerFunc(h.GETTariffProfileList))
	mux.Handle("HEAD /tp", http.HandlerFunc(h.HEADTariffProfileList))
	mux.Handle("POST /tp", http.HandlerFunc(h.POSTTariffProfileList))
	mux.Handle("PUT /tp", http.HandlerFunc(h.PUTTariffProfileList))
	mux.Handle("DELETE /tp/{id1}", http.HandlerFunc(h.DELETETariffProfile))
	mux.Handle("GET /tp/{id1}", http.HandlerFunc(h.GETTariffProfile))
	mux.Handle("HEAD /tp/{id1}", http.HandlerFunc(h.HEADTariffProfile))
	mux.Handle("POST /tp/{id1}", http.HandlerFunc(h.POSTTariffProfile))
	mux.Handle("PUT /tp/{id1}", http.HandlerFunc(h.PUTTariffProfile))
	mux.Handle("DELETE /tp/{id1}/rc", http.HandlerFunc(h.DELETERateComponentList))
	mux.Handle("GET /tp/{id1}/rc", http.HandlerFunc(h.GETRateComponentList))
	mux.Handle("HEAD /tp/{id1}/rc", http.HandlerFunc(h.HEADRateComponentList))
	mux.Handle("POST /tp/{id1}/rc", http.HandlerFunc(h.POSTRateComponentList))
	mux.Handle("PUT /tp/{id1}/rc", http.HandlerFunc(h.PUTRateComponentList))
	mux.Handle("DELETE /tp/{id1}/rc/{id2}", http.HandlerFunc(h.DELETERateComponent))
	mux.Handle("GET /tp/{id1}/rc/{id2}", http.HandlerFunc(h.GETRateComponent))
	mux.Handle("HEAD /tp/{id1}/rc/{id2}", http.HandlerFunc(h.HEADRateComponent))
	mux.Handle("POST /tp/{id1}/rc/{id2}", http.HandlerFunc(h.POSTRateComponent))
	mux.Handle("PUT /tp/{id1}/rc/{id2}", http.HandlerFunc(h.PUTRateComponent))
	mux.Handle("DELETE /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.DELETEActiveTimeTariffIntervalList))
	mux.Handle("GET /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.GETActiveTimeTariffIntervalList))
	mux.Handle("HEAD /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.HEADActiveTimeTariffIntervalList))
	mux.Handle("POST /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.POSTActiveTimeTariffIntervalList))
	mux.Handle("PUT /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.PUTActiveTimeTariffIntervalList))
	mux.Handle("DELETE /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.DELETETimeTariffIntervalList))
	mux.Handle("GET /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.GETTimeTariffIntervalList))
	mux.Handle("HEAD /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.HEADTimeTariffIntervalList))
	mux.Handle("POST /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.POSTTimeTariffIntervalList))
	mux.Handle("PUT /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.PUTTimeTariffIntervalList))
	mux.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.DELETETimeTariffInterval))
	mux.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.GETTimeTariffInterval))
	mux.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.HEADTimeTariffInterval))
	mux.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.POSTTimeTariffInterval))
	mux.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.PUTTimeTariffInterval))
	mux.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.DELETEConsumptionTariffIntervalList))
	mux.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.GETConsumptionTariffIntervalList))
	mux.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.HEADConsumptionTariffIntervalList))
	mux.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.POSTConsumptionTariffIntervalList))
	mux.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.PUTConsumptionTariffIntervalList))
	mux.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.DELETEConsumptionTariffInterval))
	mux.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.GETConsumptionTariffInterval))
	mux.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.HEADConsumptionTariffInterval))
	mux.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.POSTConsumptionTariffInterval))
	mux.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.PUTConsumptionTariffInterval))
}

// TestTPRouteRegistration verifies all Tp routes are correctly registered.
func TestTPRouteRegistration(t *testing.T) {
	// Create test mux with Tp routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerTPRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/tp_routes.json")
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

// TestTPPathParameters verifies path parameter extraction.
func TestTPPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/tp": {},
		"/tp/{id1}": {"id1"},
		"/tp/{id1}/rc": {"id1"},
		"/tp/{id1}/rc/{id2}": {"id1", "id2"},
		"/tp/{id1}/rc/{id2}/acttti": {"id1", "id2"},
		"/tp/{id1}/rc/{id2}/tti": {"id1", "id2"},
		"/tp/{id1}/rc/{id2}/tti/{id3}": {"id1", "id2", "id3"},
		"/tp/{id1}/rc/{id2}/tti/{id3}/cti": {"id1", "id2", "id3"},
		"/tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}": {"id1", "id2", "id3", "id4"},
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

// TestTPHTTPMethods verifies HTTP method support.
func TestTPHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerTPRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Load test data for verification
	routeData, err := ioutil.ReadFile("../../test/testdata/tp_routes.json")
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
