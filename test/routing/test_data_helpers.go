package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// LoadTestData loads route expectations from a JSON test data file.
func LoadTestData(fileName string) ([]RouteExpectation, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read test data file: %w", err)
	}

	var testData struct {
		ServiceName string `json:"service_name"`
		WADLFile    string `json:"wadl_file"`
		TotalRoutes int    `json:"total_routes"`
		UniquePaths int    `json:"unique_paths"`
		Routes      []struct {
			Method string `json:"Method"`
			Path   string `json:"Path"`
		} `json:"routes"`
	}

	err = json.Unmarshal(data, &testData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test data: %w", err)
	}

	expectations := make([]RouteExpectation, 0, len(testData.Routes))
	for _, r := range testData.Routes {
		expectations = append(expectations, RouteExpectation{
			Method: r.Method,
			Path:   r.Path,
		})
	}

	return expectations, nil
}

// LoadTestDataByServiceName loads test data using service name.
// Handles special naming conventions.
func LoadTestDataByServiceName(serviceName string) ([]RouteExpectation, error) {
	// Map service names to test data file names
	fileMap := map[string]string{
		"Bill":          "bill_routes.json",
		"BRS":           "brs_routes.json",
		"DCAP":          "dcap_routes.json",
		"DERP":          "derp_routes.json",
		"DR":            "dr_routes.json",
		"EDevice":       "edev_routes.json",
		"File":          "file_routes.json",
		"MUP":           "mup_routes.json",
		"Messaging":     "msg_routes.json",
		"Notify":        "ntfy_routes.json",
		"PPY":           "ppy_routes.json",
		"SDevice":       "sdev_routes.json",
		"TariffProfile": "tp_routes.json",
		"TimeOfUse":     "tm_routes.json",
		"UPT":           "upt_routes.json",
		"rsps":          "rsps_routes.json",
	}

	fileName, exists := fileMap[serviceName]
	if !exists {
		return nil, fmt.Errorf("unknown service: %s", serviceName)
	}

	filePath := filepath.Join("test", "testdata", fileName)
	return LoadTestData(filePath)
}

// GroupRoutesByMethod groups routes by their HTTP method.
func GroupRoutesByMethod(routes []RouteExpectation) map[string][]RouteExpectation {
	grouped := make(map[string][]RouteExpectation)

	for _, route := range routes {
		grouped[route.Method] = append(grouped[route.Method], route)
	}

	return grouped
}

// GetMethodsForPath returns all methods for a specific path.
func GetMethodsForPath(routes []RouteExpectation, path string) []string {
	methods := make([]string, 0)
	seen := make(map[string]bool)

	for _, route := range routes {
		if route.Path == path && !seen[route.Method] {
			methods = append(methods, route.Method)
			seen[route.Method] = true
		}
	}

	return methods
}

// GetPathsForMethod returns all paths for a specific method.
func GetPathsForMethod(routes []RouteExpectation, method string) []string {
	paths := make([]string, 0)
	seen := make(map[string]bool)

	for _, route := range routes {
		if route.Method == method && !seen[route.Path] {
			paths = append(paths, route.Path)
			seen[route.Path] = true
		}
	}

	return paths
}

// HTTPMethodTestCase represents a test case for HTTP method verification.
type HTTPMethodTestCase struct {
	Method             string
	Path               string
	ExpectedSupported  bool
	Description        string
}

// GenerateHTTPMethodTestCases creates test cases from route expectations.
func GenerateHTTPMethodTestCases(routes []RouteExpectation) []HTTPMethodTestCase {
	testCases := make([]HTTPMethodTestCase, 0, len(routes))

	for _, route := range routes {
		testCases = append(testCases, HTTPMethodTestCase{
			Method:            route.Method,
			Path:              route.Path,
			ExpectedSupported: true,
			Description:       fmt.Sprintf("%s %s should be supported", route.Method, route.Path),
		})
	}

	return testCases
}

// UnsupportedMethodTestCase represents a test for unsupported methods.
type UnsupportedMethodTestCase struct {
	Method     string
	Path       string
	Description string
}

// GenerateUnsupportedMethodTestCases creates negative test cases.
func GenerateUnsupportedMethodTestCases(routes []RouteExpectation) []UnsupportedMethodTestCase {
	testCases := make([]UnsupportedMethodTestCase, 0)
	pathSet := make(map[string]bool)

	// Collect all paths
	for _, route := range routes {
		pathSet[route.Path] = true
	}

	// Generate negative tests for common unsupported methods
	unsupportedMethods := []string{"PATCH", "OPTIONS", "TRACE", "CONNECT"}

	for path := range pathSet {
		for _, method := range unsupportedMethods {
			testCases = append(testCases, UnsupportedMethodTestCase{
				Method:      method,
				Path:        path,
				Description: fmt.Sprintf("%s %s should not be supported (or return 405)", method, path),
			})
		}
	}

	return testCases
}
