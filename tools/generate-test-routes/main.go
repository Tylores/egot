package main

import (
"encoding/json"
"flag"
"fmt"
"io/ioutil"
"os"
"path/filepath"
"strings"

"github.com/Tylores/egot/test/routing"
)

type RouteTestData struct {
ServiceName string                  `json:"service_name"`
WADLFile    string                  `json:"wadl_file"`
TotalRoutes int                     `json:"total_routes"`
UniquePaths int                     `json:"unique_paths"`
Routes      []routing.ExtractedRoute `json:"routes"`
}

func main() {
serviceFlag := flag.String("service", "", "Service name (e.g., Bill, EDevice). Empty means all.")
flag.Parse()

wadlDir := "./wadl"
outputDir := "./test/testdata"

// Ensure output directory exists
os.MkdirAll(outputDir, 0755)

extractor := routing.NewWADLExtractor(wadlDir)

var services []string
if *serviceFlag != "" {
services = []string{*serviceFlag}
} else {
// Get all WADL files
files, err := ioutil.ReadDir(wadlDir)
if err != nil {
fmt.Fprintf(os.Stderr, "Error reading wadl directory: %v\n", err)
os.Exit(1)
}

for _, file := range files {
if strings.HasSuffix(file.Name(), ".wadl") {
serviceName := strings.TrimSuffix(file.Name(), ".wadl")
// Capitalize service name
serviceName = strings.ToUpper(serviceName[:1]) + serviceName[1:]
services = append(services, serviceName)
}
}
}

if len(services) == 0 {
fmt.Fprintf(os.Stderr, "No services found\n")
os.Exit(1)
}

fmt.Printf("Generating test data for %d service(s)...\n", len(services))

totalProcessed := 0
totalRoutes := 0

for _, service := range services {
// Extract routes
routes, err := extractor.ExtractRoutesFromServiceWADL(service)
if err != nil {
fmt.Fprintf(os.Stderr, "Warning: Failed to extract routes for %s: %v\n", service, err)
continue
}

// Deduplicate and sort
unique := routing.DeduplicateRoutes(routes)
sorted := routing.SortRoutes(unique)

// Convert to expectations
_ = routing.ExtractRoutesToExpectations(sorted) // for reference

// Group by path
grouped := routing.GroupRoutesByPath(unique)

// Create test data
testData := RouteTestData{
ServiceName: service,
WADLFile:    strings.ToLower(service) + ".wadl",
TotalRoutes: len(sorted),
UniquePaths: len(grouped),
Routes:      sorted,
}

// Save to JSON
outputFile := filepath.Join(outputDir, strings.ToLower(service)+"_routes.json")
data, err := json.MarshalIndent(testData, "", "  ")
if err != nil {
fmt.Fprintf(os.Stderr, "Error marshaling JSON for %s: %v\n", service, err)
continue
}

err = ioutil.WriteFile(outputFile, data, 0644)
if err != nil {
fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", outputFile, err)
continue
}

fmt.Printf("✓ %s: %d routes (%d paths) -> %s\n", 
service, testData.TotalRoutes, testData.UniquePaths, outputFile)

totalProcessed++
totalRoutes += testData.TotalRoutes
}

fmt.Printf("\n✅ Generated test data for %d service(s)\n", totalProcessed)
fmt.Printf("   Total routes: %d\n", totalRoutes)
}
