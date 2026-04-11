package scaffold

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// sep2Application is the root element of a SEP 2 WADL file.
type sep2Application struct {
	XMLName   xml.Name    `xml:"application"`
	Resources sep2ResRoot `xml:"resources"`
	Attrs     []xml.Attr  `xml:",any,attr"`
}

type sep2ResRoot struct {
	Resources []sep2Resource `xml:"resource"`
}

type sep2Resource struct {
	ID      string       `xml:"id,attr"`
	Methods []sep2Method `xml:"method"`
	Attrs   []xml.Attr   `xml:",any,attr"`
}

// UnmarshalXML extracts the namespaced wx:samplePath attribute.
func (r *sep2Resource) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		if a.Name.Local == "samplePath" {
			r.Attrs = append(r.Attrs, a)
		}
	}
	type Alias sep2Resource
	return d.DecodeElement((*Alias)(r), &start)
}

func (r *sep2Resource) samplePath() string {
	for _, a := range r.Attrs {
		if a.Name.Local == "samplePath" {
			return a.Value
		}
	}
	return ""
}

type sep2Method struct {
	ID       string         `xml:"id,attr"`
	Name     string         `xml:"name,attr"`
	Response []sep2Response `xml:"response"`
	Attrs    []xml.Attr     `xml:",any,attr"`
}

func (m *sep2Method) mode() string {
	for _, a := range m.Attrs {
		if a.Name.Local == "mode" {
			return a.Value
		}
	}
	return "M"
}

type sep2Response struct {
	Reps []sep2Representation `xml:"representation"`
}

type sep2Representation struct {
	Element string `xml:"element,attr"`
}

// internal service description built from the parsed WADL

type serviceSpec struct {
	MaxEntities int
	Resources   []resourceSpec
}

type resourceSpec struct {
	Path    string
	Name    string // resource id (e.g. "BillingReadingSetList")
	Methods []methodSpec
}

type methodSpec struct {
	FuncName     string   // from method id (e.g. "GETDeviceCapability")
	HTTPMethod   string   // GET, POST, PUT, DELETE, HEAD
	Mode         string   // M, D, E, O
	ResponseType string   // sep type name (e.g. "DeviceCapability"), empty if none
	PathParams   []string // param names extracted from samplePath
}

// Generator generates microservice scaffolds from WADL specifications.
type Generator struct {
	spec        serviceSpec
	serviceName string
}

// NewGenerator creates a new generator from a SEP 2 WADL file.
// The service name is derived from the WADL filename (without extension).
func NewGenerator(wadlPath string) (*Generator, error) {
	data, err := os.ReadFile(wadlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read WADL file: %w", err)
	}

	app := &sep2Application{}
	if err = xml.Unmarshal(data, app); err != nil {
		return nil, fmt.Errorf("failed to parse WADL: %w", err)
	}

	spec, err := buildSpec(app)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(wadlPath)
	serviceName := strings.TrimSuffix(base, filepath.Ext(base))

	return &Generator{spec: spec, serviceName: serviceName}, nil
}

// buildSpec converts a parsed SEP 2 application into the internal service spec.
func buildSpec(app *sep2Application) (serviceSpec, error) {
	spec := serviceSpec{MaxEntities: 100}

	for _, a := range app.Attrs {
		switch a.Name.Local {
		case "max_entities":
			v, _ := strconv.Atoi(a.Value)
			if v > 0 {
				spec.MaxEntities = v
			}
		}
	}

	for _, res := range app.Resources.Resources {
		path := res.samplePath()
		if path == "" {
			continue
		}

		pathParams := extractPathParams(path)
		rspec := resourceSpec{
			Path: path,
			Name: res.ID,
		}

		for _, m := range res.Methods {
			if m.ID == "" {
				continue
			}
			mspec := methodSpec{
				FuncName:   m.ID,
				HTTPMethod: strings.ToUpper(m.Name),
				Mode:       m.mode(),
				PathParams: pathParams,
			}
			// Extract response type from first representation with an element attr
			for _, resp := range m.Response {
				for _, rep := range resp.Reps {
					if rep.Element != "" {
						parts := strings.SplitN(rep.Element, ":", 2)
						mspec.ResponseType = parts[len(parts)-1]
						break
					}
				}
				if mspec.ResponseType != "" {
					break
				}
			}
			rspec.Methods = append(rspec.Methods, mspec)
		}

		if len(rspec.Methods) > 0 {
			spec.Resources = append(spec.Resources, rspec)
		}
	}

	return spec, nil
}

// extractPathParams returns param names from a path like "/brs/{id1}/br/{id2}".
func extractPathParams(path string) []string {
	var params []string
	for _, seg := range strings.Split(path, "/") {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			params = append(params, seg[1:len(seg)-1])
		}
	}
	return params
}

// ServiceName returns the service name derived from the WADL filename.
func (g *Generator) ServiceName() string {
	return g.serviceName
}

// Generate creates the complete service scaffold in the current directory.
func (g *Generator) Generate() error {
	serviceName := g.serviceName
	outputDir := "."

	dirs := []string{
		filepath.Join(outputDir, fmt.Sprintf("cmd/%s", serviceName)),
		filepath.Join(outputDir, fmt.Sprintf("internal/%s/handler", serviceName)),
		filepath.Join(outputDir, fmt.Sprintf("internal/%s/repository/memory", serviceName)),
		filepath.Join(outputDir, fmt.Sprintf("internal/%s/server", serviceName)),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	for _, fn := range []func(string, string) error{
		g.generateMain,
		g.generateHandler,
		g.generateRepository,
		g.generateRepositoryError,
		g.generateServer,
		g.updateRoutes,
	} {
		if err := fn(outputDir, serviceName); err != nil {
			return err
		}
	}

	return nil
}

// toCamelCase converts "device-manager" → "DeviceManager".
func toCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// toConstantFormat converts "device-manager" → "DEVICE_MANAGER".
func toConstantFormat(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}
