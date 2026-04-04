package routing

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// WADLResource represents a WADL resource element.
type WADLResource struct {
	ID         string          `xml:"id,attr"`
	Path       string          `xml:"path,attr"`
	SamplePath string          `xml:"samplePath,attr"`
	Methods    []WADLMethod    `xml:"method"`
	Resources  []WADLResource  `xml:"resource"`
}

// WADLMethod represents a WADL method element.
type WADLMethod struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// WADLApplication is the root WADL element.
type WADLApplication struct {
	XMLName xml.Name       `xml:"application"`
	Resources []WADLResource `xml:"resources>resource"`
}

// ExtractedRoute represents an extracted route from WADL.
type ExtractedRoute struct {
	Method string // HTTP method (GET, POST, PUT, DELETE)
	Path   string // URL path
}

// WADLExtractor extracts routes from WADL files.
type WADLExtractor struct {
	wadlDir string
}

// NewWADLExtractor creates a new WADLExtractor for a wadl directory.
func NewWADLExtractor(wadlDir string) *WADLExtractor {
	return &WADLExtractor{wadlDir: wadlDir}
}

// ExtractRoutesFromFile extracts routes from a single WADL file.
func (we *WADLExtractor) ExtractRoutesFromFile(fileName string) ([]ExtractedRoute, error) {
	filePath := filepath.Join(we.wadlDir, fileName)
	
	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read WADL file %s: %w", filePath, err)
	}

	// Parse XML
	var app WADLApplication
	err = xml.Unmarshal(data, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WADL file %s: %w", filePath, err)
	}

	// Extract routes
	routes := make([]ExtractedRoute, 0)
	for _, resource := range app.Resources {
		extracted := we.extractFromResource(resource, "")
		routes = append(routes, extracted...)
	}

	return routes, nil
}

// extractFromResource recursively extracts routes from a resource and its children.
func (we *WADLExtractor) extractFromResource(resource WADLResource, parentPath string) []ExtractedRoute {
	routes := make([]ExtractedRoute, 0)

	// Build current path
	currentPath := parentPath
	if resource.SamplePath != "" {
		// Use samplePath if available (SEP2 WADL format)
		currentPath = resource.SamplePath
	} else if resource.Path != "" {
		// Use path attribute
		if strings.HasPrefix(resource.Path, "/") {
			currentPath = resource.Path
		} else {
			currentPath = parentPath + "/" + resource.Path
		}
	}

	// Extract methods for this resource
	for _, method := range resource.Methods {
		if method.Name != "" && currentPath != "" {
			routes = append(routes, ExtractedRoute{
				Method: strings.ToUpper(method.Name),
				Path:   currentPath,
			})
		}
	}

	// Recursively process child resources
	for _, childResource := range resource.Resources {
		childRoutes := we.extractFromResource(childResource, currentPath)
		routes = append(routes, childRoutes...)
	}

	return routes
}

// ExtractRoutesFromServiceWADL extracts routes for a specific service.
// Service name is used to find the WADL file (e.g., "Bill" -> "bill.wadl").
func (we *WADLExtractor) ExtractRoutesFromServiceWADL(serviceName string) ([]ExtractedRoute, error) {
	// Build WADL file name from service name
	wadlFileName := strings.ToLower(serviceName) + ".wadl"
	
	return we.ExtractRoutesFromFile(wadlFileName)
}

// ExtractRoutesToExpectations converts extracted routes to RouteExpectation format.
func ExtractRoutesToExpectations(routes []ExtractedRoute) []RouteExpectation {
	expectations := make([]RouteExpectation, 0, len(routes))
	
	for _, route := range routes {
		expectations = append(expectations, RouteExpectation{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	
	return expectations
}

// DeduplicateRoutes removes duplicate routes.
func DeduplicateRoutes(routes []ExtractedRoute) []ExtractedRoute {
	seen := make(map[string]bool)
	result := make([]ExtractedRoute, 0)
	
	for _, route := range routes {
		key := fmt.Sprintf("%s %s", route.Method, route.Path)
		if !seen[key] {
			seen[key] = true
			result = append(result, route)
		}
	}
	
	return result
}

// GroupRoutesByPath groups routes by their path.
// Returns map: path -> []methods
func GroupRoutesByPath(routes []ExtractedRoute) map[string][]string {
	grouped := make(map[string][]string)
	
	for _, route := range routes {
		if _, exists := grouped[route.Path]; !exists {
			grouped[route.Path] = make([]string, 0)
		}
		grouped[route.Path] = append(grouped[route.Path], route.Method)
	}
	
	return grouped
}

// GetRouteCount returns the count of unique routes.
func GetRouteCount(routes []ExtractedRoute) int {
	unique := make(map[string]bool)
	
	for _, route := range routes {
		key := fmt.Sprintf("%s %s", route.Method, route.Path)
		unique[key] = true
	}
	
	return len(unique)
}

// GetPathCount returns the count of unique paths (regardless of method).
func GetPathCount(routes []ExtractedRoute) int {
	unique := make(map[string]bool)
	
	for _, route := range routes {
		unique[route.Path] = true
	}
	
	return len(unique)
}

// ValidateRoutes checks for common issues in extracted routes.
// Returns (valid bool, errors []string)
func ValidateRoutes(routes []ExtractedRoute) (bool, []string) {
	errors := make([]string, 0)
	
	for i, route := range routes {
		// Check method
		if route.Method == "" {
			errors = append(errors, fmt.Sprintf("Route %d: empty method", i))
		}
		
		validMethods := map[string]bool{
			"GET": true, "POST": true, "PUT": true, "DELETE": true,
			"HEAD": true, "PATCH": true, "OPTIONS": true,
		}
		if !validMethods[route.Method] {
			errors = append(errors, fmt.Sprintf("Route %d: invalid method '%s'", i, route.Method))
		}
		
		// Check path
		if route.Path == "" {
			errors = append(errors, fmt.Sprintf("Route %d: empty path", i))
		}
		
		if !strings.HasPrefix(route.Path, "/") {
			errors = append(errors, fmt.Sprintf("Route %d: path must start with '/', got '%s'", i, route.Path))
		}
	}
	
	return len(errors) == 0, errors
}

// SortRoutes sorts routes for consistent output.
// Sorts by: method (alphabetically), then path (alphabetically).
func SortRoutes(routes []ExtractedRoute) []ExtractedRoute {
	// Simple bubble sort for demonstration
	// In production, use sort.Slice
	result := make([]ExtractedRoute, len(routes))
	copy(result, routes)
	
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			a := fmt.Sprintf("%s %s", result[i].Method, result[i].Path)
			b := fmt.Sprintf("%s %s", result[j].Method, result[j].Path)
			if a > b {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	return result
}

// FormatRoute formats a route as "METHOD /path".
func FormatRoute(route ExtractedRoute) string {
	return fmt.Sprintf("%s %s", route.Method, route.Path)
}

// FormatRoutes formats multiple routes, one per line.
func FormatRoutes(routes []ExtractedRoute) string {
	lines := make([]string, 0, len(routes))
	
	for _, route := range routes {
		lines = append(lines, FormatRoute(route))
	}
	
	return strings.Join(lines, "\n")
}
