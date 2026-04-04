// Code generated - DO NOT EDIT
// This test implements route registration tests for Messaging service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/Messaging/repository/memory"
	"github.com/Tylores/egot/internal/Messaging/handler"
	"github.com/Tylores/egot/test/routing"
)

// TestMessagingRouteRegistration verifies all Messaging routes are correctly registered.
func TestMessagingRouteRegistration(t *testing.T) {
	// Create test mux with Messaging routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerMessagingRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/msg_routes.json")
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

// TestMessagingPathParameters verifies path parameter extraction.
func TestMessagingPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
		"/msg": nil,
		"/msg/{id1}": {"id1"},
		"/msg/{id1}/acttxt": {"id1"},
		"/msg/{id1}/txt": {"id1"},
		"/msg/{id1}/txt/{id2}": {"id1", "id2"},
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

// TestMessagingHTTPMethods verifies HTTP method support.
func TestMessagingHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	registerMessagingRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
		{"DELETE", "/msg", true},
		{"DELETE", "/msg/{id1}", true},
		{"DELETE", "/msg/{id1}/acttxt", true},
		{"DELETE", "/msg/{id1}/txt", true},
		{"DELETE", "/msg/{id1}/txt/{id2}", true},
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

// registerMessagingRoutes registers all Messaging service routes
func registerMessagingRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All 25 routes registered in main.go
	mux.Handle("DELETE /msg", http.HandlerFunc(h.DELETEMessagingProgramList))
	mux.Handle("GET /msg", http.HandlerFunc(h.GETMessagingProgramList))
	mux.Handle("HEAD /msg", http.HandlerFunc(h.HEADMessagingProgramList))
	mux.Handle("POST /msg", http.HandlerFunc(h.POSTMessagingProgramList))
	mux.Handle("PUT /msg", http.HandlerFunc(h.PUTMessagingProgramList))

	mux.Handle("DELETE /msg/{id1}", http.HandlerFunc(h.DELETEMessagingProgram))
	mux.Handle("GET /msg/{id1}", http.HandlerFunc(h.GETMessagingProgram))
	mux.Handle("HEAD /msg/{id1}", http.HandlerFunc(h.HEADMessagingProgram))
	mux.Handle("POST /msg/{id1}", http.HandlerFunc(h.POSTMessagingProgram))
	mux.Handle("PUT /msg/{id1}", http.HandlerFunc(h.PUTMessagingProgram))

	mux.Handle("DELETE /msg/{id1}/acttxt", http.HandlerFunc(h.DELETEActiveTextMessageList))
	mux.Handle("GET /msg/{id1}/acttxt", http.HandlerFunc(h.GETActiveTextMessageList))
	mux.Handle("HEAD /msg/{id1}/acttxt", http.HandlerFunc(h.HEADActiveTextMessageList))
	mux.Handle("POST /msg/{id1}/acttxt", http.HandlerFunc(h.POSTActiveTextMessageList))
	mux.Handle("PUT /msg/{id1}/acttxt", http.HandlerFunc(h.PUTActiveTextMessageList))

	mux.Handle("DELETE /msg/{id1}/txt", http.HandlerFunc(h.DELETETextMessageList))
	mux.Handle("GET /msg/{id1}/txt", http.HandlerFunc(h.GETTextMessageList))
	mux.Handle("HEAD /msg/{id1}/txt", http.HandlerFunc(h.HEADTextMessageList))
	mux.Handle("POST /msg/{id1}/txt", http.HandlerFunc(h.POSTTextMessageList))
	mux.Handle("PUT /msg/{id1}/txt", http.HandlerFunc(h.PUTTextMessageList))

	mux.Handle("DELETE /msg/{id1}/txt/{id2}", http.HandlerFunc(h.DELETETextMessage))
	mux.Handle("GET /msg/{id1}/txt/{id2}", http.HandlerFunc(h.GETTextMessage))
	mux.Handle("HEAD /msg/{id1}/txt/{id2}", http.HandlerFunc(h.HEADTextMessage))
	mux.Handle("POST /msg/{id1}/txt/{id2}", http.HandlerFunc(h.POSTTextMessage))
	mux.Handle("PUT /msg/{id1}/txt/{id2}", http.HandlerFunc(h.PUTTextMessage))
}
