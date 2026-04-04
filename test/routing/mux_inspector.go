package routing

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// MuxInspector provides utilities for inspecting HTTP mux route registrations.
type MuxInspector struct {
	mux    *http.ServeMux
	routes map[string]*RouteInfo
}

// RouteInfo contains information about a registered route.
type RouteInfo struct {
	Method   string
	Path     string
	Handler  http.Handler
	Matched  bool // Whether a test request matched this route
}

// NewMuxInspector creates a new MuxInspector for analyzing a mux.
func NewMuxInspector(mux *http.ServeMux) *MuxInspector {
	return &MuxInspector{
		mux:    mux,
		routes: make(map[string]*RouteInfo),
	}
}

// DiscoverRoutes attempts to discover registered routes by testing common patterns.
// This is a heuristic approach since ServeMux doesn't expose registered routes directly.
// Returns a list of discovered routes.
func (mi *MuxInspector) DiscoverRoutes(testPaths []string) []RouteInfo {
	discovered := make([]RouteInfo, 0)

	for _, path := range testPaths {
		for _, method := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"} {
			matched := mi.TestRoute(method, path)
			if matched {
				info := RouteInfo{
					Method:  method,
					Path:    path,
					Matched: true,
				}
				key := fmt.Sprintf("%s %s", method, path)
				mi.routes[key] = &info
				discovered = append(discovered, info)
			}
		}
	}

	return discovered
}

// TestRoute attempts to match a route and returns whether it matched.
func (mi *MuxInspector) TestRoute(method, path string) bool {
	// Create a test request
	req := httptest.NewRequest(method, path, nil)

	// Create a response recorder to avoid side effects
	recorder := httptest.NewRecorder()

	// Make request and see if it's handled vs. 404
	mi.mux.ServeHTTP(recorder, req)

	// If we didn't get a 404, the route was matched
	matched := recorder.Code != http.StatusNotFound

	return matched
}

// findHandlerForRequest attempts to identify which handler will handle a request.
// This is a helper for route discovery.
func (mi *MuxInspector) findHandlerForRequest(r *http.Request) http.Handler {
	// This is difficult without access to mux internals
	// We'll rely on the TestRoute method instead
	return nil
}

// VerifyRoutesRegistered checks if all expected routes are registered.
// Returns (allRegistered bool, missing []string, registered []string)
func (mi *MuxInspector) VerifyRoutesRegistered(expectations []RouteExpectation, testPaths []string) (bool, []string, []string) {
	// Discover what's actually registered
	discovered := mi.DiscoverRoutes(testPaths)

	// Build map of discovered routes
	discoveredMap := make(map[string]bool)
	for _, route := range discovered {
		key := fmt.Sprintf("%s %s", route.Method, route.Path)
		discoveredMap[key] = true
	}

	// Build map of expected routes
	expectedMap := make(map[string]bool)
	for _, exp := range expectations {
		key := fmt.Sprintf("%s %s", exp.Method, exp.Path)
		expectedMap[key] = true
	}

	// Find missing routes
	missing := make([]string, 0)
	for _, exp := range expectations {
		key := fmt.Sprintf("%s %s", exp.Method, exp.Path)
		if !discoveredMap[key] {
			missing = append(missing, key)
		}
	}

	// Get registered routes
	registered := make([]string, 0, len(discoveredMap))
	for route := range discoveredMap {
		registered = append(registered, route)
	}

	return len(missing) == 0, missing, registered
}

// GetRegisteredRoutes returns all discovered registered routes.
func (mi *MuxInspector) GetRegisteredRoutes() []RouteInfo {
	routes := make([]RouteInfo, 0, len(mi.routes))
	for _, info := range mi.routes {
		routes = append(routes, *info)
	}
	return routes
}

// GetRegisteredCount returns the number of discovered routes.
func (mi *MuxInspector) GetRegisteredCount() int {
	return len(mi.routes)
}

// RequestLogger provides detailed logging of HTTP requests for debugging.
type RequestLogger struct {
	output io.Writer
}

// NewRequestLogger creates a new RequestLogger.
func NewRequestLogger(output io.Writer) *RequestLogger {
	return &RequestLogger{output: output}
}

// LogRequest logs details about an HTTP request.
func (rl *RequestLogger) LogRequest(req *http.Request, response *http.Response) {
	if rl.output == nil {
		return
	}

	fmt.Fprintf(rl.output, "Request: %s %s\n", req.Method, req.URL.Path)
	fmt.Fprintf(rl.output, "  Host: %s\n", req.Host)
	fmt.Fprintf(rl.output, "  Path: %s\n", req.URL.Path)

	if response != nil {
		fmt.Fprintf(rl.output, "Response: %d\n", response.StatusCode)
	}
	fmt.Fprintf(rl.output, "\n")
}

// TestRequest performs an HTTP request and returns matched status.
func (mi *MuxInspector) TestRequest(method, path string) (int, error) {
	req := httptest.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	mi.mux.ServeHTTP(recorder, req)

	return recorder.Code, nil
}

// TestPathWithParameters tests a path with actual parameter values.
// Example: "/bill/{id1}" becomes "/bill/123"
func (mi *MuxInspector) TestPathWithParameters(method, pathPattern string, paramValues map[string]string) (int, error) {
	path := pathPattern
	for paramName, paramValue := range paramValues {
		placeholder := fmt.Sprintf("{%s}", paramName)
		path = strings.ReplaceAll(path, placeholder, paramValue)
	}

	return mi.TestRequest(method, path)
}

// AssertRouteMatchesMethod verifies that a route is registered for a specific method.
// Returns true if the route matches the method.
func (mi *MuxInspector) AssertRouteMatchesMethod(method, path string) bool {
	status, _ := mi.TestRequest(method, path)
	return status != http.StatusNotFound
}

// AssertUnsupportedMethod verifies that a method is NOT supported (should return 405 or similar).
// Returns true if the method is not supported.
func (mi *MuxInspector) AssertUnsupportedMethod(method, path string) bool {
	status, _ := mi.TestRequest(method, path)
	// Could be 405 Method Not Allowed or 404 Not Found (depending on how mux handles it)
	return status != http.StatusOK && status != http.StatusCreated && status != http.StatusNoContent
}

// ComparisonResult holds the result of comparing expected vs actual routes.
type ComparisonResult struct {
	Total      int
	Registered int
	MissingCount int
	ExtraCount int
	Routes     []string
	MissingRoutes []string
	ExtraRoutes   []string
}

// CompareRoutes compares expected routes vs discovered routes.
func (mi *MuxInspector) CompareRoutes(expectations []RouteExpectation, testPaths []string) ComparisonResult {
	result := ComparisonResult{
		Total:   len(expectations),
		Routes:  make([]string, 0),
		MissingRoutes: make([]string, 0),
		ExtraRoutes:   make([]string, 0),
	}

	// Discover actual routes
	discovered := mi.DiscoverRoutes(testPaths)
	result.Registered = len(discovered)

	// Build maps for comparison
	expectedMap := make(map[string]bool)
	for _, exp := range expectations {
		key := fmt.Sprintf("%s %s", exp.Method, exp.Path)
		expectedMap[key] = true
	}

	discoveredMap := make(map[string]bool)
	for _, route := range discovered {
		key := fmt.Sprintf("%s %s", route.Method, route.Path)
		discoveredMap[key] = true
		result.Routes = append(result.Routes, key)
	}

	// Find missing and extra
	for key := range expectedMap {
		if !discoveredMap[key] {
			result.MissingRoutes = append(result.MissingRoutes, key)
		}
	}

	for key := range discoveredMap {
		if !expectedMap[key] {
			result.ExtraRoutes = append(result.ExtraRoutes, key)
		}
	}

	result.MissingCount = len(result.MissingRoutes)
	result.ExtraCount = len(result.ExtraRoutes)

	return result
}
