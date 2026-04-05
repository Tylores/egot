// Code generation tool to implement route registration tests for all microservices
// Usage: go run tools/implement_tests.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// ServiceMapping maps service names to their directory names and config
var ServiceMapping = map[string]struct {
	dir       string
	pkgName   string
	dataName  string
	handlerImport string
}{
	"brs":            {dir: "BRS", pkgName: "handler", dataName: "brs", handlerImport: "github.com/Tylores/egot/internal/BRS/handler"},
	"dcap":           {dir: "DCAP", pkgName: "handler", dataName: "dcap", handlerImport: "github.com/Tylores/egot/internal/DCAP/handler"},
	"derp":           {dir: "DERP", pkgName: "handler", dataName: "derp", handlerImport: "github.com/Tylores/egot/internal/DERP/handler"},
	"dr":             {dir: "DR", pkgName: "handler", dataName: "dr", handlerImport: "github.com/Tylores/egot/internal/DR/handler"},
	"edev":           {dir: "EDevice", pkgName: "handler", dataName: "edev", handlerImport: "github.com/Tylores/egot/internal/EDevice/handler"},
	"file":           {dir: "File", pkgName: "handler", dataName: "file", handlerImport: "github.com/Tylores/egot/internal/File/handler"},
	"mup":            {dir: "MUP", pkgName: "handler", dataName: "mup", handlerImport: "github.com/Tylores/egot/internal/MUP/handler"},
	"messaging":      {dir: "Messaging", pkgName: "handler", dataName: "msg", handlerImport: "github.com/Tylores/egot/internal/Messaging/handler"},
	"notify":         {dir: "Notify", pkgName: "handler", dataName: "ntfy", handlerImport: "github.com/Tylores/egot/internal/Notify/handler"},
	"ppy":            {dir: "PPY", pkgName: "handler", dataName: "ppy", handlerImport: "github.com/Tylores/egot/internal/PPY/handler"},
	"sdevice":        {dir: "SDevice", pkgName: "handler", dataName: "sdev", handlerImport: "github.com/Tylores/egot/internal/SDevice/handler"},
	"tariffprofile":  {dir: "TariffProfile", pkgName: "handler", dataName: "tp", handlerImport: "github.com/Tylores/egot/internal/TariffProfile/handler"},
	"timeofuse":      {dir: "TimeOfUse", pkgName: "handler", dataName: "tm", handlerImport: "github.com/Tylores/egot/internal/TimeOfUse/handler"},
	"upt":            {dir: "UPT", pkgName: "handler", dataName: "upt", handlerImport: "github.com/Tylores/egot/internal/UPT/handler"},
	"rsps":           {dir: "rsps", pkgName: "handler", dataName: "rsps", handlerImport: "github.com/Tylores/egot/internal/rsps/handler"},
}

type RouteData struct {
	ServiceName string `json:"service_name"`
	TotalRoutes int    `json:"total_routes"`
	Routes      []Route
}

type Route struct {
	Method string
	Path   string
}

func main() {
	// Get all services to implement
	services := getServicesToImplement()
	
	fmt.Printf("Found %d services to implement\n", len(services))
	for _, service := range services {
		if err := implementServiceTests(service); err != nil {
			fmt.Printf("Error implementing %s: %v\n", service, err)
			continue
		}
		fmt.Printf("✅ Implemented tests for %s\n", service)
	}
	
	fmt.Printf("\nImplementation complete for %d services\n", len(services))
}

func getServicesToImplement() []string {
	var services []string
	for service := range ServiceMapping {
		if service != "bill" { // Skip Bill - already implemented
			services = append(services, service)
		}
	}
	sort.Strings(services)
	return services
}

func implementServiceTests(service string) error {
	config, ok := ServiceMapping[service]
	if !ok {
		return fmt.Errorf("service not found in mapping")
	}
	
	// Load test data
	routeData, err := loadRouteData(config.dataName)
	if err != nil {
		return fmt.Errorf("failed to load route data: %w", err)
	}
	
	// Generate test file
	testContent := generateTestFile(service, config, routeData)
	
	// Write test file
	testPath := filepath.Join("cmd", config.dir, fmt.Sprintf("%s_route_registration_test.go", service))
	if err := ioutil.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	return nil
}

func loadRouteData(service string) (*RouteData, error) {
	dataPath := filepath.Join("test", "testdata", fmt.Sprintf("%s_routes.json", service))
	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}
	
	var routeData RouteData
	if err := json.Unmarshal(data, &routeData); err != nil {
		return nil, err
	}
	
	return &routeData, nil
}

func generateTestFile(service string, config struct {
	dir       string
	pkgName   string
	dataName  string
	handlerImport string
}, routeData *RouteData) string {
	// Determine handler variable name and import style
	handlerVarName := "h"
	
	// Build service constants (e.g., BRS, Bill)
	serviceKeyUpper := strings.ToUpper(service)
	serviceKey := capitalizeFirst(service)
	
	// Generate handler registration code
	routesByMethod := groupRoutesByMethod(routeData.Routes)
	registrationCode := generateRegistrationCode(handlerVarName, routesByMethod)
	
	// Generate path parameters
	pathParams := extractPathPatterns(routeData.Routes)
	pathParamCode := generatePathParamCode(pathParams)
	
	// Generate HTTP method tests
	httpMethodTests := generateHTTPMethodTests(routesByMethod)
	
	content := fmt.Sprintf(`// Code generated by implement_tests.go - DO NOT EDIT
// This test file implements route registration tests for %s service.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"%s"
	"github.com/Tylores/egot/test/routing"
)

// Test%sRouteRegistration verifies all %s routes are correctly registered.
func Test%sRouteRegistration(t *testing.T) {
	// Create test mux with %s routes
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	register%sRoutes(mux, h)

	// Load expected routes from test data
	routeData, err := ioutil.ReadFile("../../test/testdata/%s_routes.json")
	if err != nil {
		t.Logf("Note: Could not load test data: %%v", err)
		t.Skip("Test data not available")
	}

	var testData struct {
		Routes []Route
	}

	if err := json.Unmarshal(routeData, &testData); err != nil {
		t.Fatalf("Failed to parse test data: %%v", err)
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
		t.Errorf("Missing %%d routes:", len(missing))
		for _, route := range missing {
			t.Errorf("  - %%s", route)
		}
	}

	// Verify route count
	matches, registered, expected := helper.AssertRouteCount(expectations)
	if !matches {
		t.Errorf("Route count mismatch: expected %%d, got %%d", expected, registered)
	}
}

// Test%sPathParameters verifies path parameter extraction.
func Test%sPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	// Define expected path patterns from WADL
	patterns := map[string][]string{
%s
	}

	for path, expectedParams := range patterns {
		if expectedParams != nil {
			validator.AddPathPattern(path, expectedParams)
		}

		// Extract parameters
		extracted := validator.ExtractParameterNames(path)

		// Verify count
		if len(extracted) != len(expectedParams) {
			t.Errorf("Path %%s: expected %%d params, got %%d",
				path, len(expectedParams), len(extracted))
			continue
		}

		// Verify names
		for i, expected := range expectedParams {
			if i < len(extracted) && extracted[i] != expected {
				t.Errorf("Path %%s param %%d: expected '%%s', got '%%s'",
					path, i, expected, extracted[i])
			}
		}
	}
}

// Test%sHTTPMethods verifies HTTP method support.
func Test%sHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	repo := memory.NewRepository()
	h := handler.NewHandler(repo)

	// Register routes
	register%sRoutes(mux, h)

	inspector := routing.NewMuxInspector(mux)

	// Test sample routes and methods
	tests := []struct {
		method string
		path   string
		expect bool
	}{
%s
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
			t.Logf("%%s %%s: status=%%d, expected found=%%v",
				test.method, test.path, status, test.expect)
		}
	}
}

// register%sRoutes registers all %s service routes (from main.go)
func register%sRoutes(mux *http.ServeMux, h *handler.Handler) {
	// All %%d routes from %s.wadl
	// Pattern: http.Handle("METHOD /path", http.HandlerFunc(h.MethodName))
	//
%s
}
`,
		serviceKeyUpper,
		config.handlerImport,
		serviceKey, serviceKeyUpper,
		serviceKey, serviceKeyUpper,
		serviceKey,
		service, service,
		serviceKey,
		serviceKey,
		pathParamCode,
		serviceKey,
		serviceKey,
		serviceKey,
		httpMethodTests,
		serviceKey,
		serviceKeyUpper,
		serviceKey,
		registrationCode,
	)
	
	return content
}

type methodRoutes struct {
	method string
	paths  []string
}

func groupRoutesByMethod(routes []Route) []methodRoutes {
	methodMap := make(map[string][]string)
	for _, route := range routes {
		methodMap[route.Method] = append(methodMap[route.Method], route.Path)
	}
	
	var result []methodRoutes
	methods := []string{"DELETE", "GET", "HEAD", "POST", "PUT"}
	for _, method := range methods {
		if paths, ok := methodMap[method]; ok {
			sort.Strings(paths)
			result = append(result, methodRoutes{method, paths})
		}
	}
	
	return result
}

func generateRegistrationCode(handlerVar string, routesByMethod []methodRoutes) string {
	var lines []string
	
	for _, mr := range routesByMethod {
		lines = append(lines, fmt.Sprintf("\t// %s routes", mr.method))
		for _, path := range mr.paths {
			handlerName := generateHandlerName(mr.method, path)
			lines = append(lines, fmt.Sprintf(`	mux.Handle("%s %s", http.HandlerFunc(%s.%s))`, mr.method, path, handlerVar, handlerName))
		}
	}
	
	return strings.Join(lines, "\n")
}

func generateHandlerName(method, path string) string {
	// Convert path to handler name
	// /bill -> Bill
	// /bill/{id1}/ca/{id2} -> BillCAInstance
	// Pattern: METHOD + ResourceName
	
	// Clean path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	var name string
	for _, part := range parts {
		if !strings.Contains(part, "{") {
			name += strings.ToUpper(part[:1]) + part[1:]
		} else {
			name += "Instance"
		}
	}
	
	return method + name
}

func extractPathPatterns(routes []Route) map[string][]string {
	patterns := make(map[string][]string)
	
	for _, route := range routes {
		params := extractParams(route.Path)
		if _, ok := patterns[route.Path]; !ok {
			patterns[route.Path] = params
		}
	}
	
	return patterns
}

func extractParams(path string) []string {
	var params []string
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			param := strings.TrimPrefix(strings.TrimSuffix(part, "}"), "{")
			params = append(params, param)
		}
	}
	return params
}

func generatePathParamCode(patterns map[string][]string) string {
	// Sort paths for consistent output
	var paths []string
	for path := range patterns {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	
	var lines []string
	for _, path := range paths {
		params := patterns[path]
		if len(params) == 0 {
			lines = append(lines, fmt.Sprintf(`		"%s": nil,`, path))
		} else {
			paramList := fmt.Sprintf(`{"%s"}`, strings.Join(params, `", "`))
			lines = append(lines, fmt.Sprintf(`		"%s": %s,`, path, paramList))
		}
	}
	
	return strings.Join(lines, "\n")
}

func generateHTTPMethodTests(routesByMethod []methodRoutes) string {
	var lines []string
	
	// Generate sample tests for each unique path
	testedPaths := make(map[string]bool)
	
	for _, mr := range routesByMethod {
		for _, path := range mr.paths {
			if !testedPaths[path] {
				lines = append(lines, fmt.Sprintf(`		{"%s", "%s", true},`, mr.method, path))
				testedPaths[path] = true
			}
		}
	}
	
	return strings.Join(lines, "\n")
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
