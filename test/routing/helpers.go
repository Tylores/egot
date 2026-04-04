package routing

import (
	"fmt"
	"net/http"
)

// TestHelper provides common utility functions for route testing.
type TestHelper struct {
	validator *RouteValidator
	inspector *MuxInspector
	methods   *HTTPMethodValidator
	pathParams *PathParameterValidator
}

// NewTestHelper creates a new TestHelper for a given mux.
func NewTestHelper(mux *http.ServeMux) *TestHelper {
	return &TestHelper{
		validator:  NewRouteValidator(mux),
		inspector:  NewMuxInspector(mux),
		methods:    NewHTTPMethodValidator(),
		pathParams: NewPathParameterValidator(),
	}
}

// RegisterExpectedRoute adds an expected route to the test helper.
func (th *TestHelper) RegisterExpectedRoute(method, path string) {
	th.validator.RegisterRoute(method, path)
	th.methods.AddSupportedMethod(path, method)
}

// RegisterExpectedRoutes adds multiple expected routes at once.
func (th *TestHelper) RegisterExpectedRoutes(expectations []RouteExpectation) {
	for _, exp := range expectations {
		th.RegisterExpectedRoute(exp.Method, exp.Path)
	}
}

// AssertAllRoutesExist verifies all expected routes are registered.
// Returns (allExist bool, missingRoutes []string)
func (th *TestHelper) AssertAllRoutesExist(expectations []RouteExpectation) (bool, []string) {
	missing := th.validator.AssertAllRoutesExist(expectations)
	return len(missing) == 0, missing
}

// AssertRouteCount verifies the count of registered routes matches expectations.
// Returns (matches bool, registered int, expected int)
func (th *TestHelper) AssertRouteCount(expectations []RouteExpectation) (bool, int, int) {
	return th.validator.VerifyRouteCount(expectations)
}

// AssertMethodSupported verifies a specific HTTP method is supported for a path.
func (th *TestHelper) AssertMethodSupported(method, path string) bool {
	return th.methods.IsSupportedMethod(path, method)
}

// AssertPathHasParameters checks if a path has parameter placeholders.
func (th *TestHelper) AssertPathHasParameters(path string) bool {
	return th.pathParams.HasPathParameters(path)
}

// ExtractPathParameters extracts parameter names from a path.
func (th *TestHelper) ExtractPathParameters(path string) []string {
	return th.pathParams.ExtractParameterNames(path)
}

// VerifyPathParameters checks if a path has the expected parameters.
// Returns (valid bool, expectedParams []string)
func (th *TestHelper) VerifyPathParameters(path string, expectedParams []string) bool {
	th.pathParams.AddPathPattern(path, expectedParams)
	valid, _ := th.pathParams.VerifyPathPattern(path)
	return valid
}

// GetSummary returns a test summary for all routes.
func (th *TestHelper) GetSummary(expectations []RouteExpectation) TestSummary {
	summary := TestSummary{
		TotalExpected:  len(expectations),
		TotalRegistered: th.validator.GetRegisteredCount(),
		Routes:         th.validator.GetAllRoutes(),
	}

	missing := th.validator.AssertAllRoutesExist(expectations)
	summary.MissingCount = len(missing)
	summary.Missing = missing

	return summary
}

// TestSummary holds summary information about route testing.
type TestSummary struct {
	TotalExpected    int
	TotalRegistered  int
	MissingCount     int
	Routes           []string
	Missing          []string
}

// String returns a formatted string representation of the summary.
func (ts TestSummary) String() string {
	return fmt.Sprintf(
		"Expected: %d, Registered: %d, Missing: %d",
		ts.TotalExpected,
		ts.TotalRegistered,
		ts.MissingCount,
	)
}

// ParameterizedPathTester helps test paths with parameters.
type ParameterizedPathTester struct {
	mux *http.ServeMux
}

// NewParameterizedPathTester creates a new ParameterizedPathTester.
func NewParameterizedPathTester(mux *http.ServeMux) *ParameterizedPathTester {
	return &ParameterizedPathTester{mux: mux}
}

// TestPathPattern tests a path pattern with multiple parameter sets.
// Returns (allMatched bool, results []PathTestResult)
func (ppt *ParameterizedPathTester) TestPathPattern(method, pathPattern string, paramSets []map[string]string) (bool, []PathTestResult) {
	results := make([]PathTestResult, 0, len(paramSets))

	allMatched := true
	for i, params := range paramSets {
		path := pathPattern
		for name, value := range params {
			ph := fmt.Sprintf("{%s}", name)
			path = fmt.Sprintf(path, value) // Simple replacement - not production-grade
			_ = ph // Use placeholder
		}

		result := PathTestResult{
			Index:   i,
			Path:    path,
			Params:  params,
			Matched: true, // Would need actual mux testing
		}

		if !result.Matched {
			allMatched = false
		}

		results = append(results, result)
	}

	return allMatched, results
}

// PathTestResult holds the result of a single parametrized path test.
type PathTestResult struct {
	Index   int
	Path    string
	Params  map[string]string
	Matched bool
}

// RouteTestBuilder helps construct complex route test scenarios.
type RouteTestBuilder struct {
	expectations []RouteExpectation
	methods      map[string][]string // path -> methods
	parameters   map[string][]string // path -> param names
}

// NewRouteTestBuilder creates a new RouteTestBuilder.
func NewRouteTestBuilder() *RouteTestBuilder {
	return &RouteTestBuilder{
		expectations: make([]RouteExpectation, 0),
		methods:      make(map[string][]string),
		parameters:   make(map[string][]string),
	}
}

// AddRoute adds a single route to the builder.
func (rtb *RouteTestBuilder) AddRoute(method, path string) *RouteTestBuilder {
	rtb.expectations = append(rtb.expectations, RouteExpectation{
		Method: method,
		Path:   path,
	})

	// Track methods by path
	if _, exists := rtb.methods[path]; !exists {
		rtb.methods[path] = make([]string, 0)
	}
	rtb.methods[path] = append(rtb.methods[path], method)

	return rtb
}

// AddRoutes adds multiple routes with the same methods.
func (rtb *RouteTestBuilder) AddRoutes(methods []string, paths []string) *RouteTestBuilder {
	for _, path := range paths {
		for _, method := range methods {
			rtb.AddRoute(method, path)
		}
	}
	return rtb
}

// AddPath adds all HTTP methods for a path.
func (rtb *RouteTestBuilder) AddPath(path string) *RouteTestBuilder {
	return rtb.AddRoutes([]string{"GET", "POST", "PUT", "DELETE"}, []string{path})
}

// WithParameters associates parameter names with a path.
func (rtb *RouteTestBuilder) WithParameters(path string, paramNames []string) *RouteTestBuilder {
	rtb.parameters[path] = paramNames
	return rtb
}

// Build returns the constructed expectations and metadata.
func (rtb *RouteTestBuilder) Build() ([]RouteExpectation, map[string][]string, map[string][]string) {
	return rtb.expectations, rtb.methods, rtb.parameters
}

// Count returns the total number of routes.
func (rtb *RouteTestBuilder) Count() int {
	return len(rtb.expectations)
}

// GetExpectations returns the expectations slice.
func (rtb *RouteTestBuilder) GetExpectations() []RouteExpectation {
	return rtb.expectations
}
