package routing

import (
	"fmt"
	"net/http"
	"strings"
)

// RouteExpectation represents an expected HTTP route.
type RouteExpectation struct {
	Method string // HTTP method (GET, POST, PUT, DELETE, etc.)
	Path   string // URL path (e.g., "/bill", "/bill/{id1}")
}

// RouteValidator provides methods for validating route registration.
type RouteValidator struct {
	registered map[string]bool // key: "METHOD /path"
	routes     []string        // all registered routes
}

// NewRouteValidator creates a new RouteValidator for a given mux.
func NewRouteValidator(mux *http.ServeMux) *RouteValidator {
	return &RouteValidator{
		registered: make(map[string]bool),
		routes:     make([]string, 0),
	}
}

// RegisterRoute marks a route as expected to be registered.
// This is typically called during test setup to populate expected routes.
func (rv *RouteValidator) RegisterRoute(method, path string) {
	key := fmt.Sprintf("%s %s", strings.ToUpper(method), path)
	rv.registered[key] = true
	rv.routes = append(rv.routes, key)
}

// AssertRouteExists checks if a route is registered.
// Returns true if the route exists, false otherwise.
func (rv *RouteValidator) AssertRouteExists(method, path string) bool {
	key := fmt.Sprintf("%s %s", strings.ToUpper(method), path)
	return rv.registered[key]
}

// AssertAllRoutesExist checks if all expected routes are registered.
// Returns a list of missing routes, or empty slice if all exist.
func (rv *RouteValidator) AssertAllRoutesExist(expectations []RouteExpectation) []string {
	missing := make([]string, 0)
	for _, exp := range expectations {
		if !rv.AssertRouteExists(exp.Method, exp.Path) {
			missing = append(missing, fmt.Sprintf("%s %s", exp.Method, exp.Path))
		}
	}
	return missing
}

// GetRegisteredCount returns the total number of registered routes.
func (rv *RouteValidator) GetRegisteredCount() int {
	return len(rv.registered)
}

// GetExpectedCount returns the total number of expected routes from last verification.
func (rv *RouteValidator) GetExpectedCount() int {
	return len(rv.routes)
}

// VerifyRouteCount ensures the number of registered routes matches expectations.
// Returns (matches bool, registered int, expected int)
func (rv *RouteValidator) VerifyRouteCount(expectations []RouteExpectation) (bool, int, int) {
	return len(rv.registered) == len(expectations), len(rv.registered), len(expectations)
}

// GetAllRoutes returns a list of all registered routes.
func (rv *RouteValidator) GetAllRoutes() []string {
	routes := make([]string, 0, len(rv.registered))
	for route := range rv.registered {
		routes = append(routes, route)
	}
	return routes
}

// HasDuplicates checks if there are any duplicate routes.
// Returns true if duplicates found.
func (rv *RouteValidator) HasDuplicates() bool {
	seen := make(map[string]bool)
	for route := range rv.registered {
		if seen[route] {
			return true
		}
		seen[route] = true
	}
	return false
}

// HTTPMethodValidator validates HTTP method support for routes.
type HTTPMethodValidator struct {
	supportedMethods map[string]map[string]bool // path -> method -> supported
}

// NewHTTPMethodValidator creates a new HTTPMethodValidator.
func NewHTTPMethodValidator() *HTTPMethodValidator {
	return &HTTPMethodValidator{
		supportedMethods: make(map[string]map[string]bool),
	}
}

// AddSupportedMethod marks a method as supported for a path.
func (hmv *HTTPMethodValidator) AddSupportedMethod(path, method string) {
	if hmv.supportedMethods[path] == nil {
		hmv.supportedMethods[path] = make(map[string]bool)
	}
	hmv.supportedMethods[path][strings.ToUpper(method)] = true
}

// IsSupportedMethod checks if a method is supported for a path.
func (hmv *HTTPMethodValidator) IsSupportedMethod(path, method string) bool {
	if hmv.supportedMethods[path] == nil {
		return false
	}
	return hmv.supportedMethods[path][strings.ToUpper(method)]
}

// GetSupportedMethods returns all supported methods for a path.
func (hmv *HTTPMethodValidator) GetSupportedMethods(path string) []string {
	methods := make([]string, 0)
	if hmv.supportedMethods[path] == nil {
		return methods
	}
	for method := range hmv.supportedMethods[path] {
		methods = append(methods, method)
	}
	return methods
}

// PathParameterValidator validates path parameters.
type PathParameterValidator struct {
	patterns map[string][]string // path -> list of parameter names
}

// NewPathParameterValidator creates a new PathParameterValidator.
func NewPathParameterValidator() *PathParameterValidator {
	return &PathParameterValidator{
		patterns: make(map[string][]string),
	}
}

// AddPathPattern registers a path pattern with its parameter names.
// Example: "/bill/{id1}/ca/{id2}" with ["id1", "id2"]
func (ppv *PathParameterValidator) AddPathPattern(path string, paramNames []string) {
	ppv.patterns[path] = paramNames
}

// HasPathParameters checks if a path has parameters.
func (ppv *PathParameterValidator) HasPathParameters(path string) bool {
	return strings.Contains(path, "{") && strings.Contains(path, "}")
}

// ExtractParameterNames extracts parameter names from a path.
// Example: "/bill/{id1}/ca/{id2}" returns ["id1", "id2"]
func (ppv *PathParameterValidator) ExtractParameterNames(path string) []string {
	params := make([]string, 0)
	start := 0
	for {
		idx := strings.Index(path[start:], "{")
		if idx == -1 {
			break
		}
		start += idx + 1
		end := strings.Index(path[start:], "}")
		if end == -1 {
			break
		}
		paramName := path[start : start+end]
		params = append(params, paramName)
		start += end + 1
	}
	return params
}

// VerifyPathPattern checks if a path matches the expected pattern.
func (ppv *PathParameterValidator) VerifyPathPattern(path string) (bool, []string) {
	expected, exists := ppv.patterns[path]
	if !exists {
		return true, []string{} // No pattern defined, assume correct
	}

	actual := ppv.ExtractParameterNames(path)
	if len(actual) != len(expected) {
		return false, expected
	}

	for i, param := range actual {
		if param != expected[i] {
			return false, expected
		}
	}
	return true, expected
}

// GetPathParameters returns registered parameters for a path.
func (ppv *PathParameterValidator) GetPathParameters(path string) []string {
	if params, exists := ppv.patterns[path]; exists {
		return params
	}
	return ppv.ExtractParameterNames(path)
}
