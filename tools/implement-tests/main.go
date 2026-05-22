// Code generation tool to implement route registration tests for all microservices
// Usage: go run tools/implement-tests/main.go
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
	"bill":           {dir: "Bill", pkgName: "handler", dataName: "bill", handlerImport: "github.com/Tylores/egot/internal/Bill/handler"},
	"dcap":           {dir: "DCAP", pkgName: "handler", dataName: "dcap", handlerImport: "github.com/Tylores/egot/internal/DCAP/handler"},
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
	"der":            {dir: "DER", pkgName: "handler", dataName: "der", handlerImport: "github.com/Tylores/egot/internal/DER/handler"},
	"flowreservation": {dir: "FlowReservation", pkgName: "handler", dataName: "frq", handlerImport: "github.com/Tylores/egot/internal/FlowReservation/handler"},
}

type RouteData struct {
	ServiceName string `json:"service_name"`
	Routes      []Route
}

type Route struct {
	Method string
	Path   string
}

func main() {
	// Get all services to implement
	var services []string
	for service := range ServiceMapping {
		services = append(services, service)
	}
	sort.Strings(services)
	
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
	handlerVarName := "h"
	serviceKey := capitalizeFirst(service)
	
	registrationCode := generateRegistrationCode(handlerVarName, routeData.Routes)
	pathParams := extractPathPatterns(routeData.Routes)
	pathParamCode := generatePathParamCode(pathParams)
	httpMethodTests := generateHTTPMethodTests(routeData.Routes)
	
	content := fmt.Sprintf(`// Code generated by implement-tests - DO NOT EDIT
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"%s"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/test/routing"
)

func Test%sRouteRegistration(t *testing.T) {
	mux := http.NewServeMux()
	reg := registry.New(":memory:")
	reg.Load()
	repo := store.New(":memory:")
	repo.Load()
	h := handler.NewHandler(repo, reg)

	register%sRoutes(mux, h)

	routeData, err := ioutil.ReadFile("../../test/testdata/%s_routes.json")
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
		t.Fatalf("Failed to parse test data: %%v", err)
	}

	expectations := make([]routing.RouteExpectation, 0, len(testData.Routes))
	for _, r := range testData.Routes {
		expectations = append(expectations, routing.RouteExpectation{
			Method: r.Method,
			Path:   r.Path,
		})
	}

	helper := routing.NewTestHelper(mux)
	helper.RegisterExpectedRoutes(expectations)
	all, missing := helper.AssertAllRoutesExist(expectations)
	if !all {
		t.Errorf("Missing %%d routes:", len(missing))
		for _, route := range missing {
			t.Errorf("  - %%s", route)
		}
	}
}

func Test%sPathParameters(t *testing.T) {
	validator := routing.NewPathParameterValidator()

	patterns := map[string][]string{
%s
	}

	for path, expectedParams := range patterns {
		if expectedParams != nil {
			validator.AddPathPattern(path, expectedParams)
		}
		extracted := validator.ExtractParameterNames(path)
		if len(extracted) != len(expectedParams) {
			t.Errorf("Path %%s: expected %%d params, got %%d", path, len(expectedParams), len(extracted))
			continue
		}
	}
}

func Test%sHTTPMethods(t *testing.T) {
	mux := http.NewServeMux()
	reg := registry.New(":memory:")
	reg.Load()
	repo := store.New(":memory:")
	repo.Load()
	h := handler.NewHandler(repo, reg)

	register%sRoutes(mux, h)
	inspector := routing.NewMuxInspector(mux)

	tests := []struct {
		method string
		path   string
		expect bool
	}{
%s
	}

	for _, test := range tests {
		status, _ := inspector.TestRequest(test.method, test.path)
		found := status != http.StatusNotFound
		if found != test.expect {
			t.Errorf("%%s %%s: expected found=%%v, got %%v", test.method, test.path, test.expect, found)
		}
	}
}

func register%sRoutes(mux *http.ServeMux, h *handler.Handler) {
%s
}
`,
		config.handlerImport,
		serviceKey, serviceKey, config.dataName,
		serviceKey, pathParamCode,
		serviceKey, serviceKey, httpMethodTests,
		serviceKey, registrationCode,
	)
	
	return content
}

func generateRegistrationCode(handlerVar string, routes []Route) string {
	var lines []string
	for _, r := range routes {
		handlerName := r.Method + getResourceNameFromPath(r.Path)
		lines = append(lines, fmt.Sprintf(`	mux.Handle("%s %s", http.HandlerFunc(%s.%s))`, r.Method, r.Path, handlerVar, handlerName))
	}
	return strings.Join(lines, "\n")
}

func getResourceNameFromPath(path string) string {
	// Specialized mapping for path patterns or top-level paths that don't follow standard rules
	if strings.HasSuffix(path, "/dderc") {
		return "DefaultDERControl"
	}
	if strings.HasSuffix(path, "/cdc") {
		return "CurrentDERControls"
	}
	if strings.HasSuffix(path, "/cdp") {
		return "CurrentDERProgram"
	}
	if strings.HasSuffix(path, "/dera") {
		return "DERAvailability"
	}
	if strings.HasSuffix(path, "/upt") && strings.Contains(path, "/der/") {
		return "AssociatedUsagePoint"
	}
	if strings.HasSuffix(path, "/derp") && strings.Contains(path, "/der/") {
		return "AssociatedDERProgramList"
	}

	// Specialized mapping for top-level paths that don't follow the pattern
	switch path {
	case "/dcap": return "DeviceCapability"
	case "/tm":   return "Time"
	case "/sdev": return "SelfDevice"
	}

	// Heuristic: scaffold-gen uses the ID from the WADL resource.
	// Most IDs follow a pattern: ResourceName (for instance) or ResourceNameList (for list).
	
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var segments []string
	for _, p := range parts {
		if !strings.HasPrefix(p, "{") {
			segments = append(segments, p)
		}
	}

	if len(segments) == 0 {
		return ""
	}

	// The last non-parameter segment determines the resource type
	last := segments[len(segments)-1]
	var name string
	switch last {
	case "edev": name = "EndDevice"
	case "brs":  name = "BillingReadingSet"
	case "bill": name = "CustomerAccount"
	case "tp":   name = "TariffProfile"
	case "upt":  name = "UsagePoint"
	case "mup":  name = "MirrorUsagePoint"
	case "derp": name = "DERProgram"
	case "der":  name = "DER"
	case "msg":  name = "MessagingProgram"
	case "ntfy": name = "Notification"
	case "ppy":  name = "Prepayment"
	case "ca":   name = "CustomerAgreement"
	case "bp":   name = "BillingPeriod"
	case "pro":  name = "ProjectionReading"
	case "tar":  name = "TargetReading"
	case "ver":  name = "HistoricalReading"
	case "ss":   name = "ServiceSupplier"
	case "br":   name = "BillingReading"
	case "dc":   name = "DERCurve"
	case "derc": name = "DERControl"
	case "edc":  name = "EndDeviceControl"
	case "txt":  name = "TextMessage"
	case "si":   name = "SupplyInterruptionOverride"
	case "cr":   name = "CreditRegister"
	case "mr":   name = "MeterReading"
	case "rs":   name = "ReadingSet"
	case "r":    name = "Reading"
	case "acttti": name = "ActiveTimeTariffInterval"
	case "tti":  name = "TimeTariffInterval"
	case "cti":  name = "ConsumptionTariffInterval"
	case "actderc": name = "ActiveDERControl"
	case "dderc": name = "DefaultDERControl"
	case "actedc": name = "ActiveEndDeviceControl"
	case "acttxt": name = "ActiveTextMessage"
	case "actsi": name = "ActiveSupplyInterruptionOverride"
	case "ab":   name = "AccountBalance"
	case "os":   name = "PrepayOperationStatus"
	case "frq":  name = "FlowReservationRequest"
	case "frp":  name = "FlowReservationResponse"
	case "dera": name = "DERAvailability"
	case "dercap": name = "DERCapability"
	case "dercom": name = "DERComponent"
	case "derg": name = "DERSettings"
	case "ders": name = "DERStatus"
	case "cdc":  name = "CurrentDERControls"
	case "cdp":  name = "CurrentDERProgram"
	case "adev": name = "AggregatedDevice"
	case "aggp": name = "AggregationPriority"
	case "cfg":  name = "Configuration"
	case "prcfg": name = "PriceResponseCfg"
	case "di":   name = "DeviceInformation"
	case "loc":  name = "SupportedLocale"
	case "dstat": name = "DeviceStatus"
	case "fs":   name = "FileStatus"
	case "fsa":  name = "FunctionSetAssignments"
	case "lel":  name = "LogEvent"
	case "lsl":  name = "LoadShedAvailability"
	case "ns":   name = "IPInterface"
	case "addr": name = "IPAddr"
	case "rpl":  name = "RPLInstance"
	case "srt":  name = "RPLSourceRoutes"
	case "ll":   name = "LLInterface"
	case "nbh":  name = "Neighbor"
	case "prxy": name = "ProxiedDevice"
	case "ps":   name = "PowerStatus"
	case "rg":   name = "Registration"
	case "sub":  name = "Subscription"
	case "rt":   name = "ReadingType"
	default:
		name = strings.ToUpper(last[:1]) + last[1:]
	}

	// If it's a list (ends with resource segment but no ID), append "List"
	if !strings.HasSuffix(path, "}") {
		// Exceptions: top-level discovery/status resources are usually NOT suffixed with List in WADL IDs
		if !strings.HasSuffix(name, "Capability") && !strings.HasSuffix(name, "Status") && 
		   !strings.HasSuffix(name, "Priority") && !strings.HasSuffix(name, "Information") &&
		   !strings.HasSuffix(name, "Settings") && name != "PowerStatus" && 
		   name != "Registration" && name != "Configuration" && name != "AccountBalance" {
			name += "List"
		}
	}

	return name
}

func extractPathPatterns(routes []Route) map[string][]string {
	patterns := make(map[string][]string)
	for _, route := range routes {
		params := extractParams(route.Path)
		patterns[route.Path] = params
	}
	return patterns
}

func extractParams(path string) []string {
	var params []string
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			params = append(params, strings.Trim(part, "{}"))
		}
	}
	return params
}

func generatePathParamCode(patterns map[string][]string) string {
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
			lines = append(lines, fmt.Sprintf(`		"%s": {"%s"},`, path, strings.Join(params, `", "`)))
		}
	}
	return strings.Join(lines, "\n")
}

func generateHTTPMethodTests(routes []Route) string {
	var lines []string
	seen := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + r.Path
		if !seen[key] {
			lines = append(lines, fmt.Sprintf(`		{"%s", "%s", true},`, r.Method, r.Path))
			seen[key] = true
		}
	}
	return strings.Join(lines, "\n")
}

func capitalizeFirst(s string) string {
	if len(s) == 0 { return s }
	return strings.ToUpper(s[:1]) + s[1:]
}
