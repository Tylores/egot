package scaffold

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WADL represents the parsed Web Application Description Language specification
type WADL struct {
	XMLName     xml.Name    `xml:"application"`
	Name        string      `xml:"name,attr"`
	Port        int         `xml:"port,attr"`
	MaxEntities int         `xml:"max_entities,attr"`
	Resources   []Resource  `xml:"resources>resource"`
}

// Resource represents a REST resource in the WADL
type Resource struct {
	Path    string   `xml:"path,attr"`
	Methods []Method `xml:"methods>method"`
}

// Method represents an HTTP method on a resource
type Method struct {
	Name         string      `xml:"name,attr"`
	HTTPMethod   string      `xml:"http_method,attr"`
	RequestType  string      `xml:"request_type,attr"`
	ResponseType string      `xml:"response_type,attr"`
	PathParams   []PathParam `xml:"path_params>param"`
}

// PathParam represents a path parameter
type PathParam struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// Generator generates microservice scaffolds from WADL specifications
type Generator struct {
	wadl *WADL
}

// NewGenerator creates a new generator from a WADL file
func NewGenerator(wadlPath string) (*Generator, error) {
	data, err := os.ReadFile(wadlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read WADL file: %w", err)
	}

	wadl := &WADL{}
	err = xml.Unmarshal(data, wadl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WADL: %w", err)
	}

	if wadl.Name == "" {
		return nil, fmt.Errorf("WADL must specify application name")
	}
	if wadl.Port == 0 {
		return nil, fmt.Errorf("WADL must specify application port")
	}
	if wadl.MaxEntities == 0 {
		wadl.MaxEntities = 100
	}

	return &Generator{wadl: wadl}, nil
}

// ServiceName returns the service name from the WADL
func (g *Generator) ServiceName() string {
	return g.wadl.Name
}

// Generate creates the complete service scaffold
func (g *Generator) Generate(outputDir string) error {
	serviceName := g.wadl.Name

	// Create directory structure
	dirs := []string{
		filepath.Join(outputDir, fmt.Sprintf("cmd/%s", serviceName)),
		filepath.Join(outputDir, fmt.Sprintf("internal/%s/handler", serviceName)),
		filepath.Join(outputDir, fmt.Sprintf("internal/%s/repository/memory", serviceName)),
		filepath.Join(outputDir, fmt.Sprintf("internal/%s/server", serviceName)),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate files
	err := g.generateMain(outputDir, serviceName)
	if err != nil {
		return err
	}

	err = g.generateHandler(outputDir, serviceName)
	if err != nil {
		return err
	}

	err = g.generateRepository(outputDir, serviceName)
	if err != nil {
		return err
	}

	err = g.generateRepositoryError(outputDir, serviceName)
	if err != nil {
		return err
	}

	err = g.generateServer(outputDir, serviceName)
	if err != nil {
		return err
	}

	err = g.updateRoutes(outputDir, serviceName)
	if err != nil {
		return err
	}

	return nil
}

// Helper to convert service name to constant format (e.g., "device-manager" -> "DeviceManager")
func toCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// Helper to convert to constant format (e.g., "device-manager" -> "DEVICE_MANAGER")
func toConstantFormat(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}
