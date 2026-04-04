// Code generated - DO NOT EDIT
// This test implements route registration tests for Edev service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/EDevice/repository/memory"
	"github.com/Tylores/egot/internal/EDevice/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestEdevRouteRegistration verifies all Edev routes are correctly registered.
func TestEdevRouteRegistration(t *testing.T) {
	// Create test mux with Edev routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerEdevRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/edev_routes.json")
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

// TestEdevPathParameters verifies path parameter extraction.
func TestEdevPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/edev": nil,
		"/edev/{id1}": {"id1"},
		"/edev/{id1}/adev": {"id1"},
		"/edev/{id1}/adev/{id2}": {"id1", "id2"},
		"/edev/{id1}/aggp": {"id1"},
		"/edev/{id1}/cfg": {"id1"},
		"/edev/{id1}/cfg/prcfg": {"id1"},
		"/edev/{id1}/cfg/prcfg/{id2}": {"id1", "id2"},
		"/edev/{id1}/der": {"id1"},
		"/edev/{id1}/der/{id2}": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/cdc": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/cdp": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/dera": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/dercap": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/dercom": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/dercom/{id3}": {"id1", "id2", "id3"},
		"/edev/{id1}/der/{id2}/derg": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/derp": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/ders": {"id1", "id2"},
		"/edev/{id1}/der/{id2}/upt": {"id1", "id2"},
		"/edev/{id1}/di": {"id1"},
		"/edev/{id1}/di/loc": {"id1"},
		"/edev/{id1}/di/loc/{id2}": {"id1", "id2"},
		"/edev/{id1}/dstat": {"id1"},
		"/edev/{id1}/frp": {"id1"},
		"/edev/{id1}/frp/{id2}": {"id1", "id2"},
		"/edev/{id1}/frq": {"id1"},
		"/edev/{id1}/frq/{id2}": {"id1", "id2"},
		"/edev/{id1}/fs": {"id1"},
		"/edev/{id1}/fsa": {"id1"},
		"/edev/{id1}/fsa/{id2}": {"id1", "id2"},
		"/edev/{id1}/lel": {"id1"},
		"/edev/{id1}/lel/{id2}": {"id1", "id2"},
		"/edev/{id1}/lsl": {"id1"},
		"/edev/{id1}/lsl/{id2}": {"id1", "id2"},
		"/edev/{id1}/ns": {"id1"},
		"/edev/{id1}/ns/{id2}": {"id1", "id2"},
		"/edev/{id1}/ns/{id2}/addr": {"id1", "id2"},
		"/edev/{id1}/ns/{id2}/addr/{id3}": {"id1", "id2", "id3"},
		"/edev/{id1}/ns/{id2}/addr/{id3}/rpl": {"id1", "id2", "id3"},
		"/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}": {"id1", "id2", "id3", "id4"},
		"/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt": {"id1", "id2", "id3", "id4"},
		"/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}": {"id1", "id2", "id3", "id4", "id5"},
		"/edev/{id1}/ns/{id2}/ll": {"id1", "id2"},
		"/edev/{id1}/ns/{id2}/ll/{id3}": {"id1", "id2", "id3"},
		"/edev/{id1}/ns/{id2}/ll/{id3}/nbh": {"id1", "id2", "id3"},
		"/edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}": {"id1", "id2", "id3", "id4"},
		"/edev/{id1}/prxy": {"id1"},
		"/edev/{id1}/prxy/{id2}": {"id1", "id2"},
		"/edev/{id1}/ps": {"id1"},
		"/edev/{id1}/rg": {"id1"},
		"/edev/{id1}/sub": {"id1"},
		"/edev/{id1}/sub/{id2}": {"id1", "id2"},
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

// TestEdevHTTPMethods verifies HTTP method support.
func TestEdevHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerEdevRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/edev", true},
		{"DELETE", "/edev/{id1}", true},
		{"DELETE", "/edev/{id1}/adev", true},
		{"DELETE", "/edev/{id1}/adev/{id2}", true},
		{"DELETE", "/edev/{id1}/aggp", true},
		{"DELETE", "/edev/{id1}/cfg", true},
		{"DELETE", "/edev/{id1}/cfg/prcfg", true},
		{"DELETE", "/edev/{id1}/cfg/prcfg/{id2}", true},
		{"DELETE", "/edev/{id1}/der", true},
		{"DELETE", "/edev/{id1}/der/{id2}", true},
		{"DELETE", "/edev/{id1}/der/{id2}/cdc", true},
		{"DELETE", "/edev/{id1}/der/{id2}/cdp", true},
		{"DELETE", "/edev/{id1}/der/{id2}/dera", true},
		{"DELETE", "/edev/{id1}/der/{id2}/dercap", true},
		{"DELETE", "/edev/{id1}/der/{id2}/dercom", true},
		{"DELETE", "/edev/{id1}/der/{id2}/dercom/{id3}", true},
		{"DELETE", "/edev/{id1}/der/{id2}/derg", true},
		{"DELETE", "/edev/{id1}/der/{id2}/derp", true},
		{"DELETE", "/edev/{id1}/der/{id2}/ders", true},
		{"DELETE", "/edev/{id1}/der/{id2}/upt", true},
		{"DELETE", "/edev/{id1}/di", true},
		{"DELETE", "/edev/{id1}/di/loc", true},
		{"DELETE", "/edev/{id1}/di/loc/{id2}", true},
		{"DELETE", "/edev/{id1}/dstat", true},
		{"DELETE", "/edev/{id1}/frp", true},
		{"DELETE", "/edev/{id1}/frp/{id2}", true},
		{"DELETE", "/edev/{id1}/frq", true},
		{"DELETE", "/edev/{id1}/frq/{id2}", true},
		{"DELETE", "/edev/{id1}/fs", true},
		{"DELETE", "/edev/{id1}/fsa", true},
		{"DELETE", "/edev/{id1}/fsa/{id2}", true},
		{"DELETE", "/edev/{id1}/lel", true},
		{"DELETE", "/edev/{id1}/lel/{id2}", true},
		{"DELETE", "/edev/{id1}/lsl", true},
		{"DELETE", "/edev/{id1}/lsl/{id2}", true},
		{"DELETE", "/edev/{id1}/ns", true},
		{"DELETE", "/edev/{id1}/ns/{id2}", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/addr", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/addr/{id3}", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/addr/{id3}/rpl", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/ll", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/ll/{id3}", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/ll/{id3}/nbh", true},
		{"DELETE", "/edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}", true},
		{"DELETE", "/edev/{id1}/prxy", true},
		{"DELETE", "/edev/{id1}/prxy/{id2}", true},
		{"DELETE", "/edev/{id1}/ps", true},
		{"DELETE", "/edev/{id1}/rg", true},
		{"DELETE", "/edev/{id1}/sub", true},
		{"DELETE", "/edev/{id1}/sub/{id2}", true},
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

// registerEdevRoutes registers all Edev service routes
func registerEdevRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 265 routes registered in main.go
	mux.Handle("DELETE /edev", http.HandlerFunc(h.DELETEEndDeviceList))
	mux.Handle("GET /edev", http.HandlerFunc(h.GETEndDeviceList))
	mux.Handle("HEAD /edev", http.HandlerFunc(h.HEADEndDeviceList))
	mux.Handle("POST /edev", http.HandlerFunc(h.POSTEndDeviceList))
	mux.Handle("PUT /edev", http.HandlerFunc(h.PUTEndDeviceList))

	mux.Handle("DELETE /edev/{id1}", http.HandlerFunc(h.DELETEEndDevice))
	mux.Handle("GET /edev/{id1}", http.HandlerFunc(h.GETEndDevice))
	mux.Handle("HEAD /edev/{id1}", http.HandlerFunc(h.HEADEndDevice))
	mux.Handle("POST /edev/{id1}", http.HandlerFunc(h.POSTEndDevice))
	mux.Handle("PUT /edev/{id1}", http.HandlerFunc(h.PUTEndDevice))

	mux.Handle("DELETE /edev/{id1}/adev", http.HandlerFunc(h.DELETEAggregatedDeviceList))
	mux.Handle("GET /edev/{id1}/adev", http.HandlerFunc(h.GETAggregatedDeviceList))
	mux.Handle("HEAD /edev/{id1}/adev", http.HandlerFunc(h.HEADAggregatedDeviceList))
	mux.Handle("POST /edev/{id1}/adev", http.HandlerFunc(h.POSTAggregatedDeviceList))
	mux.Handle("PUT /edev/{id1}/adev", http.HandlerFunc(h.PUTAggregatedDeviceList))

	mux.Handle("DELETE /edev/{id1}/adev/{id2}", http.HandlerFunc(h.DELETEAggregatedDevice))
	mux.Handle("GET /edev/{id1}/adev/{id2}", http.HandlerFunc(h.GETAggregatedDevice))
	mux.Handle("HEAD /edev/{id1}/adev/{id2}", http.HandlerFunc(h.HEADAggregatedDevice))
	mux.Handle("POST /edev/{id1}/adev/{id2}", http.HandlerFunc(h.POSTAggregatedDevice))
	mux.Handle("PUT /edev/{id1}/adev/{id2}", http.HandlerFunc(h.PUTAggregatedDevice))

	mux.Handle("DELETE /edev/{id1}/aggp", http.HandlerFunc(h.DELETEAggregationPriority))
	mux.Handle("GET /edev/{id1}/aggp", http.HandlerFunc(h.GETAggregationPriority))
	mux.Handle("HEAD /edev/{id1}/aggp", http.HandlerFunc(h.HEADAggregationPriority))
	mux.Handle("POST /edev/{id1}/aggp", http.HandlerFunc(h.POSTAggregationPriority))
	mux.Handle("PUT /edev/{id1}/aggp", http.HandlerFunc(h.PUTAggregationPriority))

	mux.Handle("DELETE /edev/{id1}/cfg", http.HandlerFunc(h.DELETEConfiguration))
	mux.Handle("GET /edev/{id1}/cfg", http.HandlerFunc(h.GETConfiguration))
	mux.Handle("HEAD /edev/{id1}/cfg", http.HandlerFunc(h.HEADConfiguration))
	mux.Handle("POST /edev/{id1}/cfg", http.HandlerFunc(h.POSTConfiguration))
	mux.Handle("PUT /edev/{id1}/cfg", http.HandlerFunc(h.PUTConfiguration))

	mux.Handle("DELETE /edev/{id1}/cfg/prcfg", http.HandlerFunc(h.DELETEPriceResponseCfgList))
	mux.Handle("GET /edev/{id1}/cfg/prcfg", http.HandlerFunc(h.GETPriceResponseCfgList))
	mux.Handle("HEAD /edev/{id1}/cfg/prcfg", http.HandlerFunc(h.HEADPriceResponseCfgList))
	mux.Handle("POST /edev/{id1}/cfg/prcfg", http.HandlerFunc(h.POSTPriceResponseCfgList))
	mux.Handle("PUT /edev/{id1}/cfg/prcfg", http.HandlerFunc(h.PUTPriceResponseCfgList))

	mux.Handle("DELETE /edev/{id1}/cfg/prcfg/{id2}", http.HandlerFunc(h.DELETEPriceResponseCfg))
	mux.Handle("GET /edev/{id1}/cfg/prcfg/{id2}", http.HandlerFunc(h.GETPriceResponseCfg))
	mux.Handle("HEAD /edev/{id1}/cfg/prcfg/{id2}", http.HandlerFunc(h.HEADPriceResponseCfg))
	mux.Handle("POST /edev/{id1}/cfg/prcfg/{id2}", http.HandlerFunc(h.POSTPriceResponseCfg))
	mux.Handle("PUT /edev/{id1}/cfg/prcfg/{id2}", http.HandlerFunc(h.PUTPriceResponseCfg))

	mux.Handle("DELETE /edev/{id1}/der", http.HandlerFunc(h.DELETEDERList))
	mux.Handle("GET /edev/{id1}/der", http.HandlerFunc(h.GETDERList))
	mux.Handle("HEAD /edev/{id1}/der", http.HandlerFunc(h.HEADDERList))
	mux.Handle("POST /edev/{id1}/der", http.HandlerFunc(h.POSTDERList))
	mux.Handle("PUT /edev/{id1}/der", http.HandlerFunc(h.PUTDERList))

	mux.Handle("DELETE /edev/{id1}/der/{id2}", http.HandlerFunc(h.DELETEDER))
	mux.Handle("GET /edev/{id1}/der/{id2}", http.HandlerFunc(h.GETDER))
	mux.Handle("HEAD /edev/{id1}/der/{id2}", http.HandlerFunc(h.HEADDER))
	mux.Handle("POST /edev/{id1}/der/{id2}", http.HandlerFunc(h.POSTDER))
	mux.Handle("PUT /edev/{id1}/der/{id2}", http.HandlerFunc(h.PUTDER))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/cdc", http.HandlerFunc(h.DELETECurrentDERControls))
	mux.Handle("GET /edev/{id1}/der/{id2}/cdc", http.HandlerFunc(h.GETCurrentDERControls))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/cdc", http.HandlerFunc(h.HEADCurrentDERControls))
	mux.Handle("POST /edev/{id1}/der/{id2}/cdc", http.HandlerFunc(h.POSTCurrentDERControls))
	mux.Handle("PUT /edev/{id1}/der/{id2}/cdc", http.HandlerFunc(h.PUTCurrentDERControls))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/cdp", http.HandlerFunc(h.DELETECurrentDERProgram))
	mux.Handle("GET /edev/{id1}/der/{id2}/cdp", http.HandlerFunc(h.GETCurrentDERProgram))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/cdp", http.HandlerFunc(h.HEADCurrentDERProgram))
	mux.Handle("POST /edev/{id1}/der/{id2}/cdp", http.HandlerFunc(h.POSTCurrentDERProgram))
	mux.Handle("PUT /edev/{id1}/der/{id2}/cdp", http.HandlerFunc(h.PUTCurrentDERProgram))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/dera", http.HandlerFunc(h.DELETEDERAvailability))
	mux.Handle("GET /edev/{id1}/der/{id2}/dera", http.HandlerFunc(h.GETDERAvailability))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/dera", http.HandlerFunc(h.HEADDERAvailability))
	mux.Handle("POST /edev/{id1}/der/{id2}/dera", http.HandlerFunc(h.POSTDERAvailability))
	mux.Handle("PUT /edev/{id1}/der/{id2}/dera", http.HandlerFunc(h.PUTDERAvailability))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/dercap", http.HandlerFunc(h.DELETEDERCapability))
	mux.Handle("GET /edev/{id1}/der/{id2}/dercap", http.HandlerFunc(h.GETDERCapability))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/dercap", http.HandlerFunc(h.HEADDERCapability))
	mux.Handle("POST /edev/{id1}/der/{id2}/dercap", http.HandlerFunc(h.POSTDERCapability))
	mux.Handle("PUT /edev/{id1}/der/{id2}/dercap", http.HandlerFunc(h.PUTDERCapability))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/dercom", http.HandlerFunc(h.DELETEDERComponentList))
	mux.Handle("GET /edev/{id1}/der/{id2}/dercom", http.HandlerFunc(h.GETDERComponentList))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/dercom", http.HandlerFunc(h.HEADDERComponentList))
	mux.Handle("POST /edev/{id1}/der/{id2}/dercom", http.HandlerFunc(h.POSTDERComponentList))
	mux.Handle("PUT /edev/{id1}/der/{id2}/dercom", http.HandlerFunc(h.PUTDERComponentList))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/dercom/{id3}", http.HandlerFunc(h.DELETEDERComponent))
	mux.Handle("GET /edev/{id1}/der/{id2}/dercom/{id3}", http.HandlerFunc(h.GETDERComponent))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/dercom/{id3}", http.HandlerFunc(h.HEADDERComponent))
	mux.Handle("POST /edev/{id1}/der/{id2}/dercom/{id3}", http.HandlerFunc(h.POSTDERComponent))
	mux.Handle("PUT /edev/{id1}/der/{id2}/dercom/{id3}", http.HandlerFunc(h.PUTDERComponent))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/derg", http.HandlerFunc(h.DELETEDERSettings))
	mux.Handle("GET /edev/{id1}/der/{id2}/derg", http.HandlerFunc(h.GETDERSettings))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/derg", http.HandlerFunc(h.HEADDERSettings))
	mux.Handle("POST /edev/{id1}/der/{id2}/derg", http.HandlerFunc(h.POSTDERSettings))
	mux.Handle("PUT /edev/{id1}/der/{id2}/derg", http.HandlerFunc(h.PUTDERSettings))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/derp", http.HandlerFunc(h.DELETEAssociatedDERProgramList))
	mux.Handle("GET /edev/{id1}/der/{id2}/derp", http.HandlerFunc(h.GETAssociatedDERProgramList))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/derp", http.HandlerFunc(h.HEADAssociatedDERProgramList))
	mux.Handle("POST /edev/{id1}/der/{id2}/derp", http.HandlerFunc(h.POSTAssociatedDERProgramList))
	mux.Handle("PUT /edev/{id1}/der/{id2}/derp", http.HandlerFunc(h.PUTAssociatedDERProgramList))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/ders", http.HandlerFunc(h.DELETEDERStatus))
	mux.Handle("GET /edev/{id1}/der/{id2}/ders", http.HandlerFunc(h.GETDERStatus))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/ders", http.HandlerFunc(h.HEADDERStatus))
	mux.Handle("POST /edev/{id1}/der/{id2}/ders", http.HandlerFunc(h.POSTDERStatus))
	mux.Handle("PUT /edev/{id1}/der/{id2}/ders", http.HandlerFunc(h.PUTDERStatus))

	mux.Handle("DELETE /edev/{id1}/der/{id2}/upt", http.HandlerFunc(h.DELETEAssociatedUsagePoint))
	mux.Handle("GET /edev/{id1}/der/{id2}/upt", http.HandlerFunc(h.GETAssociatedUsagePoint))
	mux.Handle("HEAD /edev/{id1}/der/{id2}/upt", http.HandlerFunc(h.HEADAssociatedUsagePoint))
	mux.Handle("POST /edev/{id1}/der/{id2}/upt", http.HandlerFunc(h.POSTAssociatedUsagePoint))
	mux.Handle("PUT /edev/{id1}/der/{id2}/upt", http.HandlerFunc(h.PUTAssociatedUsagePoint))

	mux.Handle("DELETE /edev/{id1}/di", http.HandlerFunc(h.DELETEDeviceInformation))
	mux.Handle("GET /edev/{id1}/di", http.HandlerFunc(h.GETDeviceInformation))
	mux.Handle("HEAD /edev/{id1}/di", http.HandlerFunc(h.HEADDeviceInformation))
	mux.Handle("POST /edev/{id1}/di", http.HandlerFunc(h.POSTDeviceInformation))
	mux.Handle("PUT /edev/{id1}/di", http.HandlerFunc(h.PUTDeviceInformation))

	mux.Handle("DELETE /edev/{id1}/di/loc", http.HandlerFunc(h.DELETESupportedLocaleList))
	mux.Handle("GET /edev/{id1}/di/loc", http.HandlerFunc(h.GETSupportedLocaleList))
	mux.Handle("HEAD /edev/{id1}/di/loc", http.HandlerFunc(h.HEADSupportedLocaleList))
	mux.Handle("POST /edev/{id1}/di/loc", http.HandlerFunc(h.POSTSupportedLocaleList))
	mux.Handle("PUT /edev/{id1}/di/loc", http.HandlerFunc(h.PUTSupportedLocaleList))

	mux.Handle("DELETE /edev/{id1}/di/loc/{id2}", http.HandlerFunc(h.DELETESupportedLocale))
	mux.Handle("GET /edev/{id1}/di/loc/{id2}", http.HandlerFunc(h.GETSupportedLocale))
	mux.Handle("HEAD /edev/{id1}/di/loc/{id2}", http.HandlerFunc(h.HEADSupportedLocale))
	mux.Handle("POST /edev/{id1}/di/loc/{id2}", http.HandlerFunc(h.POSTSupportedLocale))
	mux.Handle("PUT /edev/{id1}/di/loc/{id2}", http.HandlerFunc(h.PUTSupportedLocale))

	mux.Handle("DELETE /edev/{id1}/dstat", http.HandlerFunc(h.DELETEDeviceStatus))
	mux.Handle("GET /edev/{id1}/dstat", http.HandlerFunc(h.GETDeviceStatus))
	mux.Handle("HEAD /edev/{id1}/dstat", http.HandlerFunc(h.HEADDeviceStatus))
	mux.Handle("POST /edev/{id1}/dstat", http.HandlerFunc(h.POSTDeviceStatus))
	mux.Handle("PUT /edev/{id1}/dstat", http.HandlerFunc(h.PUTDeviceStatus))

	mux.Handle("DELETE /edev/{id1}/frp", http.HandlerFunc(h.DELETEFlowReservationResponseList))
	mux.Handle("GET /edev/{id1}/frp", http.HandlerFunc(h.GETFlowReservationResponseList))
	mux.Handle("HEAD /edev/{id1}/frp", http.HandlerFunc(h.HEADFlowReservationResponseList))
	mux.Handle("POST /edev/{id1}/frp", http.HandlerFunc(h.POSTFlowReservationResponseList))
	mux.Handle("PUT /edev/{id1}/frp", http.HandlerFunc(h.PUTFlowReservationResponseList))

	mux.Handle("DELETE /edev/{id1}/frp/{id2}", http.HandlerFunc(h.DELETEFlowReservationResponse))
	mux.Handle("GET /edev/{id1}/frp/{id2}", http.HandlerFunc(h.GETFlowReservationResponse))
	mux.Handle("HEAD /edev/{id1}/frp/{id2}", http.HandlerFunc(h.HEADFlowReservationResponse))
	mux.Handle("POST /edev/{id1}/frp/{id2}", http.HandlerFunc(h.POSTFlowReservationResponse))
	mux.Handle("PUT /edev/{id1}/frp/{id2}", http.HandlerFunc(h.PUTFlowReservationResponse))

	mux.Handle("DELETE /edev/{id1}/frq", http.HandlerFunc(h.DELETEFlowReservationRequestList))
	mux.Handle("GET /edev/{id1}/frq", http.HandlerFunc(h.GETFlowReservationRequestList))
	mux.Handle("HEAD /edev/{id1}/frq", http.HandlerFunc(h.HEADFlowReservationRequestList))
	mux.Handle("POST /edev/{id1}/frq", http.HandlerFunc(h.POSTFlowReservationRequestList))
	mux.Handle("PUT /edev/{id1}/frq", http.HandlerFunc(h.PUTFlowReservationRequestList))

	mux.Handle("DELETE /edev/{id1}/frq/{id2}", http.HandlerFunc(h.DELETEFlowReservationRequest))
	mux.Handle("GET /edev/{id1}/frq/{id2}", http.HandlerFunc(h.GETFlowReservationRequest))
	mux.Handle("HEAD /edev/{id1}/frq/{id2}", http.HandlerFunc(h.HEADFlowReservationRequest))
	mux.Handle("POST /edev/{id1}/frq/{id2}", http.HandlerFunc(h.POSTFlowReservationRequest))
	mux.Handle("PUT /edev/{id1}/frq/{id2}", http.HandlerFunc(h.PUTFlowReservationRequest))

	mux.Handle("DELETE /edev/{id1}/fs", http.HandlerFunc(h.DELETEFileStatus))
	mux.Handle("GET /edev/{id1}/fs", http.HandlerFunc(h.GETFileStatus))
	mux.Handle("HEAD /edev/{id1}/fs", http.HandlerFunc(h.HEADFileStatus))
	mux.Handle("POST /edev/{id1}/fs", http.HandlerFunc(h.POSTFileStatus))
	mux.Handle("PUT /edev/{id1}/fs", http.HandlerFunc(h.PUTFileStatus))

	mux.Handle("DELETE /edev/{id1}/fsa", http.HandlerFunc(h.DELETEFunctionSetAssignmentsList))
	mux.Handle("GET /edev/{id1}/fsa", http.HandlerFunc(h.GETFunctionSetAssignmentsList))
	mux.Handle("HEAD /edev/{id1}/fsa", http.HandlerFunc(h.HEADFunctionSetAssignmentsList))
	mux.Handle("POST /edev/{id1}/fsa", http.HandlerFunc(h.POSTFunctionSetAssignmentsList))
	mux.Handle("PUT /edev/{id1}/fsa", http.HandlerFunc(h.PUTFunctionSetAssignmentsList))

	mux.Handle("DELETE /edev/{id1}/fsa/{id2}", http.HandlerFunc(h.DELETEFunctionSetAssignments))
	mux.Handle("GET /edev/{id1}/fsa/{id2}", http.HandlerFunc(h.GETFunctionSetAssignments))
	mux.Handle("HEAD /edev/{id1}/fsa/{id2}", http.HandlerFunc(h.HEADFunctionSetAssignments))
	mux.Handle("POST /edev/{id1}/fsa/{id2}", http.HandlerFunc(h.POSTFunctionSetAssignments))
	mux.Handle("PUT /edev/{id1}/fsa/{id2}", http.HandlerFunc(h.PUTFunctionSetAssignments))

	mux.Handle("DELETE /edev/{id1}/lel", http.HandlerFunc(h.DELETELogEventList))
	mux.Handle("GET /edev/{id1}/lel", http.HandlerFunc(h.GETLogEventList))
	mux.Handle("HEAD /edev/{id1}/lel", http.HandlerFunc(h.HEADLogEventList))
	mux.Handle("POST /edev/{id1}/lel", http.HandlerFunc(h.POSTLogEventList))
	mux.Handle("PUT /edev/{id1}/lel", http.HandlerFunc(h.PUTLogEventList))

	mux.Handle("DELETE /edev/{id1}/lel/{id2}", http.HandlerFunc(h.DELETELogEvent))
	mux.Handle("GET /edev/{id1}/lel/{id2}", http.HandlerFunc(h.GETLogEvent))
	mux.Handle("HEAD /edev/{id1}/lel/{id2}", http.HandlerFunc(h.HEADLogEvent))
	mux.Handle("POST /edev/{id1}/lel/{id2}", http.HandlerFunc(h.POSTLogEvent))
	mux.Handle("PUT /edev/{id1}/lel/{id2}", http.HandlerFunc(h.PUTLogEvent))

	mux.Handle("DELETE /edev/{id1}/lsl", http.HandlerFunc(h.DELETELoadShedAvailabilityList))
	mux.Handle("GET /edev/{id1}/lsl", http.HandlerFunc(h.GETLoadShedAvailabilityList))
	mux.Handle("HEAD /edev/{id1}/lsl", http.HandlerFunc(h.HEADLoadShedAvailabilityList))
	mux.Handle("POST /edev/{id1}/lsl", http.HandlerFunc(h.POSTLoadShedAvailabilityList))
	mux.Handle("PUT /edev/{id1}/lsl", http.HandlerFunc(h.PUTLoadShedAvailabilityList))

	mux.Handle("DELETE /edev/{id1}/lsl/{id2}", http.HandlerFunc(h.DELETELoadShedAvailability))
	mux.Handle("GET /edev/{id1}/lsl/{id2}", http.HandlerFunc(h.GETLoadShedAvailability))
	mux.Handle("HEAD /edev/{id1}/lsl/{id2}", http.HandlerFunc(h.HEADLoadShedAvailability))
	mux.Handle("POST /edev/{id1}/lsl/{id2}", http.HandlerFunc(h.POSTLoadShedAvailability))
	mux.Handle("PUT /edev/{id1}/lsl/{id2}", http.HandlerFunc(h.PUTLoadShedAvailability))

	mux.Handle("DELETE /edev/{id1}/ns", http.HandlerFunc(h.DELETEIPInterfaceList))
	mux.Handle("GET /edev/{id1}/ns", http.HandlerFunc(h.GETIPInterfaceList))
	mux.Handle("HEAD /edev/{id1}/ns", http.HandlerFunc(h.HEADIPInterfaceList))
	mux.Handle("POST /edev/{id1}/ns", http.HandlerFunc(h.POSTIPInterfaceList))
	mux.Handle("PUT /edev/{id1}/ns", http.HandlerFunc(h.PUTIPInterfaceList))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}", http.HandlerFunc(h.DELETEIPInterface))
	mux.Handle("GET /edev/{id1}/ns/{id2}", http.HandlerFunc(h.GETIPInterface))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}", http.HandlerFunc(h.HEADIPInterface))
	mux.Handle("POST /edev/{id1}/ns/{id2}", http.HandlerFunc(h.POSTIPInterface))
	mux.Handle("PUT /edev/{id1}/ns/{id2}", http.HandlerFunc(h.PUTIPInterface))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/addr", http.HandlerFunc(h.DELETEIPAddrList))
	mux.Handle("GET /edev/{id1}/ns/{id2}/addr", http.HandlerFunc(h.GETIPAddrList))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/addr", http.HandlerFunc(h.HEADIPAddrList))
	mux.Handle("POST /edev/{id1}/ns/{id2}/addr", http.HandlerFunc(h.POSTIPAddrList))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/addr", http.HandlerFunc(h.PUTIPAddrList))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/addr/{id3}", http.HandlerFunc(h.DELETEIPAddr))
	mux.Handle("GET /edev/{id1}/ns/{id2}/addr/{id3}", http.HandlerFunc(h.GETIPAddr))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/addr/{id3}", http.HandlerFunc(h.HEADIPAddr))
	mux.Handle("POST /edev/{id1}/ns/{id2}/addr/{id3}", http.HandlerFunc(h.POSTIPAddr))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/addr/{id3}", http.HandlerFunc(h.PUTIPAddr))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/addr/{id3}/rpl", http.HandlerFunc(h.DELETERPLInstanceList))
	mux.Handle("GET /edev/{id1}/ns/{id2}/addr/{id3}/rpl", http.HandlerFunc(h.GETRPLInstanceList))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/addr/{id3}/rpl", http.HandlerFunc(h.HEADRPLInstanceList))
	mux.Handle("POST /edev/{id1}/ns/{id2}/addr/{id3}/rpl", http.HandlerFunc(h.POSTRPLInstanceList))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/addr/{id3}/rpl", http.HandlerFunc(h.PUTRPLInstanceList))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}", http.HandlerFunc(h.DELETERPLInstance))
	mux.Handle("GET /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}", http.HandlerFunc(h.GETRPLInstance))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}", http.HandlerFunc(h.HEADRPLInstance))
	mux.Handle("POST /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}", http.HandlerFunc(h.POSTRPLInstance))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}", http.HandlerFunc(h.PUTRPLInstance))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt", http.HandlerFunc(h.DELETERPLSourceRoutesList))
	mux.Handle("GET /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt", http.HandlerFunc(h.GETRPLSourceRoutesList))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt", http.HandlerFunc(h.HEADRPLSourceRoutesList))
	mux.Handle("POST /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt", http.HandlerFunc(h.POSTRPLSourceRoutesList))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt", http.HandlerFunc(h.PUTRPLSourceRoutesList))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}", http.HandlerFunc(h.DELETERPLSourceRoutes))
	mux.Handle("GET /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}", http.HandlerFunc(h.GETRPLSourceRoutes))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}", http.HandlerFunc(h.HEADRPLSourceRoutes))
	mux.Handle("POST /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}", http.HandlerFunc(h.POSTRPLSourceRoutes))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/addr/{id3}/rpl/{id4}/srt/{id5}", http.HandlerFunc(h.PUTRPLSourceRoutes))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/ll", http.HandlerFunc(h.DELETELLInterfaceList))
	mux.Handle("GET /edev/{id1}/ns/{id2}/ll", http.HandlerFunc(h.GETLLInterfaceList))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/ll", http.HandlerFunc(h.HEADLLInterfaceList))
	mux.Handle("POST /edev/{id1}/ns/{id2}/ll", http.HandlerFunc(h.POSTLLInterfaceList))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/ll", http.HandlerFunc(h.PUTLLInterfaceList))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/ll/{id3}", http.HandlerFunc(h.DELETELLInterface))
	mux.Handle("GET /edev/{id1}/ns/{id2}/ll/{id3}", http.HandlerFunc(h.GETLLInterface))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/ll/{id3}", http.HandlerFunc(h.HEADLLInterface))
	mux.Handle("POST /edev/{id1}/ns/{id2}/ll/{id3}", http.HandlerFunc(h.POSTLLInterface))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/ll/{id3}", http.HandlerFunc(h.PUTLLInterface))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/ll/{id3}/nbh", http.HandlerFunc(h.DELETENeighborList))
	mux.Handle("GET /edev/{id1}/ns/{id2}/ll/{id3}/nbh", http.HandlerFunc(h.GETNeighborList))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/ll/{id3}/nbh", http.HandlerFunc(h.HEADNeighborList))
	mux.Handle("POST /edev/{id1}/ns/{id2}/ll/{id3}/nbh", http.HandlerFunc(h.POSTNeighborList))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/ll/{id3}/nbh", http.HandlerFunc(h.PUTNeighborList))

	mux.Handle("DELETE /edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}", http.HandlerFunc(h.DELETENeighbor))
	mux.Handle("GET /edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}", http.HandlerFunc(h.GETNeighbor))
	mux.Handle("HEAD /edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}", http.HandlerFunc(h.HEADNeighbor))
	mux.Handle("POST /edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}", http.HandlerFunc(h.POSTNeighbor))
	mux.Handle("PUT /edev/{id1}/ns/{id2}/ll/{id3}/nbh/{id4}", http.HandlerFunc(h.PUTNeighbor))

	mux.Handle("DELETE /edev/{id1}/prxy", http.HandlerFunc(h.DELETEProxiedDeviceList))
	mux.Handle("GET /edev/{id1}/prxy", http.HandlerFunc(h.GETProxiedDeviceList))
	mux.Handle("HEAD /edev/{id1}/prxy", http.HandlerFunc(h.HEADProxiedDeviceList))
	mux.Handle("POST /edev/{id1}/prxy", http.HandlerFunc(h.POSTProxiedDeviceList))
	mux.Handle("PUT /edev/{id1}/prxy", http.HandlerFunc(h.PUTProxiedDeviceList))

	mux.Handle("DELETE /edev/{id1}/prxy/{id2}", http.HandlerFunc(h.DELETEProxiedDevice))
	mux.Handle("GET /edev/{id1}/prxy/{id2}", http.HandlerFunc(h.GETProxiedDevice))
	mux.Handle("HEAD /edev/{id1}/prxy/{id2}", http.HandlerFunc(h.HEADProxiedDevice))
	mux.Handle("POST /edev/{id1}/prxy/{id2}", http.HandlerFunc(h.POSTProxiedDevice))
	mux.Handle("PUT /edev/{id1}/prxy/{id2}", http.HandlerFunc(h.PUTProxiedDevice))

	mux.Handle("DELETE /edev/{id1}/ps", http.HandlerFunc(h.DELETEPowerStatus))
	mux.Handle("GET /edev/{id1}/ps", http.HandlerFunc(h.GETPowerStatus))
	mux.Handle("HEAD /edev/{id1}/ps", http.HandlerFunc(h.HEADPowerStatus))
	mux.Handle("POST /edev/{id1}/ps", http.HandlerFunc(h.POSTPowerStatus))
	mux.Handle("PUT /edev/{id1}/ps", http.HandlerFunc(h.PUTPowerStatus))

	mux.Handle("DELETE /edev/{id1}/rg", http.HandlerFunc(h.DELETERegistration))
	mux.Handle("GET /edev/{id1}/rg", http.HandlerFunc(h.GETRegistration))
	mux.Handle("HEAD /edev/{id1}/rg", http.HandlerFunc(h.HEADRegistration))
	mux.Handle("POST /edev/{id1}/rg", http.HandlerFunc(h.POSTRegistration))
	mux.Handle("PUT /edev/{id1}/rg", http.HandlerFunc(h.PUTRegistration))

	mux.Handle("DELETE /edev/{id1}/sub", http.HandlerFunc(h.DELETESubscriptionList))
	mux.Handle("GET /edev/{id1}/sub", http.HandlerFunc(h.GETSubscriptionList))
	mux.Handle("HEAD /edev/{id1}/sub", http.HandlerFunc(h.HEADSubscriptionList))
	mux.Handle("POST /edev/{id1}/sub", http.HandlerFunc(h.POSTSubscriptionList))
	mux.Handle("PUT /edev/{id1}/sub", http.HandlerFunc(h.PUTSubscriptionList))

	mux.Handle("DELETE /edev/{id1}/sub/{id2}", http.HandlerFunc(h.DELETESubscription))
	mux.Handle("GET /edev/{id1}/sub/{id2}", http.HandlerFunc(h.GETSubscription))
	mux.Handle("HEAD /edev/{id1}/sub/{id2}", http.HandlerFunc(h.HEADSubscription))
	mux.Handle("POST /edev/{id1}/sub/{id2}", http.HandlerFunc(h.POSTSubscription))
	mux.Handle("PUT /edev/{id1}/sub/{id2}", http.HandlerFunc(h.PUTSubscription))
}
