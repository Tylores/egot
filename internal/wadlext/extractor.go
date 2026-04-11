package wadlext

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// WADLApplication represents the root WADL element
type WADLApplication struct {
	XMLName   xml.Name   `xml:"application"`
	Doc       *Doc       `xml:"doc,omitempty"`
	Grammars  *Grammars  `xml:"grammars,omitempty"`
	Resources *Resources `xml:"resources,omitempty"`
	Attrs     []xml.Attr `xml:",any,attr"`
}

// Doc represents a WADL doc element
type Doc struct {
	Title string `xml:"title,attr"`
	Text  string `xml:",chardata"`
}

// Grammars represents WADL grammars section
type Grammars struct {
	Includes []Include `xml:"include"`
}

// Include represents a schema include
type Include struct {
	Href string `xml:"href,attr"`
}

// Resources represents WADL resources section
type Resources struct {
	SampleBase string     `xml:"sampleBase,attr"`
	Resources  []Resource `xml:"resource"`
	Attrs      []xml.Attr `xml:",any,attr"`
}

// Resource represents a WADL resource element
type Resource struct {
	XMLName     xml.Name      `xml:"resource"`
	ID          string        `xml:"id,attr"`
	SamplePath  string        `xml:"samplePath,attr"`
	Doc         *Doc          `xml:"doc,omitempty"`
	Methods     []Method      `xml:"method"`
	SampleParam []SampleParam `xml:"sampleParam"`
	Attrs       []xml.Attr    `xml:",any,attr"`
	RawXML      string        `xml:"-"`
}

// UnmarshalXML custom unmarshaler to handle namespaced attributes
func (r *Resource) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Resource
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	// Extract namespaced samplePath attribute
	for _, attr := range start.Attr {
		if attr.Name.Space == "urn:ieee:std:2030.5:wadlExt" && attr.Name.Local == "samplePath" {
			r.SamplePath = attr.Value
		}
	}

	err := d.DecodeElement(aux, &start)
	return err
}

// Method represents a WADL method element
type Method struct {
	XMLName  xml.Name   `xml:"method"`
	ID       string     `xml:"id,attr"`
	Name     string     `xml:"name,attr"`
	Request  *Request   `xml:"request,omitempty"`
	Response []Response `xml:"response"`
	Attrs    []xml.Attr `xml:",any,attr"`
}

// Request represents a WADL request element
type Request struct {
	XMLName xml.Name         `xml:"request"`
	Params  []Param          `xml:"param"`
	Reps    []Representation `xml:"representation"`
}

// Response represents a WADL response element
type Response struct {
	XMLName xml.Name         `xml:"response"`
	Status  string           `xml:"status,attr"`
	Params  []Param          `xml:"param"`
	Reps    []Representation `xml:"representation"`
	Attrs   []xml.Attr       `xml:",any,attr"`
}

// Param represents a WADL param element
type Param struct {
	XMLName  xml.Name   `xml:"param"`
	Name     string     `xml:"name,attr"`
	Style    string     `xml:"style,attr"`
	Type     string     `xml:"type,attr"`
	Required string     `xml:"required,attr"`
	Doc      *Doc       `xml:"doc,omitempty"`
	Attrs    []xml.Attr `xml:",any,attr"`
}

// Representation represents a WADL representation element
type Representation struct {
	XMLName   xml.Name   `xml:"representation"`
	MediaType string     `xml:"mediaType,attr"`
	Element   string     `xml:"element,attr"`
	Attrs     []xml.Attr `xml:",any,attr"`
}

// SampleParam represents a WADL sample parameter
type SampleParam struct {
	XMLName xml.Name   `xml:"sampleParam"`
	Name    string     `xml:"name,attr"`
	Style   string     `xml:"style,attr"`
	Type    string     `xml:"type,attr"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

// Extractor extracts WADL resources by path prefix
type Extractor struct {
	app    *WADLApplication
	rawXML string
}

// NewExtractor creates a new WADL extractor from a file
func NewExtractor(wadlPath string) (*Extractor, error) {
	file, err := os.Open(wadlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open WADL file: %w", err)
	}
	defer file.Close()

	return NewExtractorFromReader(file)
}

// NewExtractorFromReader creates a new WADL extractor from a reader
func NewExtractorFromReader(r io.Reader) (*Extractor, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read WADL: %w", err)
	}

	var app WADLApplication
	err = xml.Unmarshal(data, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WADL: %w", err)
	}

	return &Extractor{app: &app, rawXML: string(data)}, nil
}

// Extract extracts all resources starting with the given path prefix.
// For multiple prefixes use ExtractMany.
func (e *Extractor) Extract(pathPrefix string) *WADLApplication {
	return e.ExtractMany([]string{pathPrefix})
}

// ExtractMany extracts all resources whose SamplePath starts with any of the
// given prefixes. Resources are returned in WADL document order.
func (e *Extractor) ExtractMany(prefixes []string) *WADLApplication {
	for i, p := range prefixes {
		if !strings.HasPrefix(p, "/") {
			prefixes[i] = "/" + p
		}
	}

	extracted := &WADLApplication{
		Doc:      e.app.Doc,
		Grammars: e.app.Grammars,
		Attrs:    e.app.Attrs,
	}

	resourceMap := extractResourceBlocks(e.rawXML, prefixes)

	filtered := []Resource{}
	if e.app.Resources != nil {
		for _, res := range e.app.Resources.Resources {
			if matchesAny(res.SamplePath, prefixes) {
				res.RawXML = resourceMap[res.ID]
				filtered = append(filtered, res)
			}
		}
	}

	extracted.Resources = &Resources{
		SampleBase: "http://localhost/sep/",
		Resources:  filtered,
	}

	return extracted
}

// matchesAny returns true if path starts with any of the given prefixes.
func matchesAny(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// extractResourceBlocks extracts raw XML strings for all resources
// and returns them indexed by resource ID
func extractResourceBlocks(xmlStr string, prefixes []string) map[string]string {
	resourceMap := make(map[string]string)

	resourcePattern := regexp.MustCompile(
		`(?s)<resource\s+id="([^"]*)"[^>]*>.*?</resource>`,
	)

	samplePathPattern := regexp.MustCompile(`wx:samplePath="([^"]*)"`)

	matches := resourcePattern.FindAllStringSubmatchIndex(xmlStr, -1)
	for _, match := range matches {
		fullMatch := xmlStr[match[0]:match[1]]
		id := xmlStr[match[2]:match[3]]

		if len(prefixes) == 0 {
			resourceMap[id] = fullMatch
			continue
		}

		pathMatch := samplePathPattern.FindStringSubmatch(fullMatch)
		if len(pathMatch) > 1 && matchesAny(pathMatch[1], prefixes) {
			resourceMap[id] = fullMatch
		}
	}

	return resourceMap
}

// Save writes the extracted WADL to a file preserving raw resource XML
func (e *Extractor) Save(app *WADLApplication, outputPath string) error {
	// Extract header (everything up to and including opening <resources tag)
	resourcesStart := strings.Index(e.rawXML, "<resources")
	if resourcesStart == -1 {
		return fmt.Errorf("could not find <resources> element in original WADL")
	}

	resourcesTagEnd := strings.Index(e.rawXML[resourcesStart:], ">")
	if resourcesTagEnd == -1 {
		return fmt.Errorf("could not find end of <resources> tag")
	}
	resourcesTagEnd += resourcesStart + 1

	// Extract footer (closing resources tag and application tag)
	resourcesEnd := strings.LastIndex(e.rawXML, "</resources>")
	if resourcesEnd == -1 {
		return fmt.Errorf("could not find closing </resources> tag")
	}

	header := e.rawXML[:resourcesTagEnd]

	// Inject any attributes added via AddAttribute that are not in the original header
	header = injectApplicationAttrs(header, app.Attrs)

	footer := e.rawXML[resourcesEnd:]

	// Build resources content from raw XML blocks
	var resourcesContent strings.Builder
	for _, res := range app.Resources.Resources {
		if res.RawXML != "" {
			resourcesContent.WriteString("    ")
			resourcesContent.WriteString(res.RawXML)
			resourcesContent.WriteString("\n")
		}
	}

	// Combine
	output := header + "\n" + resourcesContent.String() + "  " + footer

	err := os.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to write WADL file: %w", err)
	}

	return nil
}

// injectApplicationAttrs inserts attributes from attrs into the <application ...> opening tag
// in header, skipping any that are already present.
func injectApplicationAttrs(header string, attrs []xml.Attr) string {
	appStart := strings.Index(header, "<application")
	if appStart == -1 {
		return header
	}
	appTagEnd := strings.Index(header[appStart:], ">")
	if appTagEnd == -1 {
		return header
	}
	appTagEnd += appStart

	appTag := header[appStart:appTagEnd]

	var toInject []string
	for _, attr := range attrs {
		// Only inject plain (non-namespaced) attributes not already present
		if attr.Name.Space == "" && !strings.Contains(appTag, attr.Name.Local+"=") {
			toInject = append(toInject, fmt.Sprintf(` %s="%s"`, attr.Name.Local, attr.Value))
		}
	}

	if len(toInject) == 0 {
		return header
	}

	return header[:appTagEnd] + strings.Join(toInject, "") + header[appTagEnd:]
}

// GetResourceCount returns the number of resources in the application
func (app *WADLApplication) GetResourceCount() int {
	if app.Resources == nil {
		return 0
	}
	return len(app.Resources.Resources)
}

// GetResourcePaths returns all resource paths in the application
func (app *WADLApplication) GetResourcePaths() []string {
	var paths []string
	if app.Resources == nil {
		return paths
	}
	for _, res := range app.Resources.Resources {
		paths = append(paths, res.SamplePath)
	}
	return paths
}

// Cluster groups resources that share a common sub-path key.
type Cluster struct {
	// Name is the distinguishing segment (or "core" for shallow paths).
	Name      string
	Resources []Resource
}

// Paths returns the samplePath of every resource in this cluster.
func (c *Cluster) Paths() []string {
	paths := make([]string, len(c.Resources))
	for i, r := range c.Resources {
		paths[i] = r.SamplePath
	}
	return paths
}

// clusterKey returns the distinguishing path segment at the given depth below
// prefix, or "core" if the path does not reach that depth.
//
// depth=1 groups by the first segment after prefix (e.g. /{id1}).
// depth=2 groups by the second segment (e.g. /{id1}/der).
func clusterKey(path, prefix string, depth int) string {
	rel := strings.TrimPrefix(path, prefix)
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" {
		return "core"
	}
	parts := strings.Split(rel, "/")
	if len(parts) < depth {
		return "core"
	}
	return parts[depth-1]
}

// SuggestClusters analyses all resources whose SamplePath starts with prefix
// and groups them by the distinguishing segment at the given depth.
// The returned slice is sorted by cluster name, with "core" first.
func (e *Extractor) SuggestClusters(prefix string, depth int) []Cluster {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	// Build a map from cluster name → resources, preserving insertion order.
	order := []string{}
	groups := map[string][]Resource{}

	if e.app.Resources != nil {
		resourceMap := extractResourceBlocks(e.rawXML, []string{prefix})
		for _, res := range e.app.Resources.Resources {
			if !strings.HasPrefix(res.SamplePath, prefix) {
				continue
			}
			res.RawXML = resourceMap[res.ID]
			key := clusterKey(res.SamplePath, prefix, depth)
			if _, exists := groups[key]; !exists {
				order = append(order, key)
			}
			groups[key] = append(groups[key], res)
		}
	}

	// Sort: "core" first, then alphabetical.
	sort.Slice(order, func(i, j int) bool {
		if order[i] == "core" {
			return true
		}
		if order[j] == "core" {
			return false
		}
		return order[i] < order[j]
	})

	clusters := make([]Cluster, 0, len(order))
	for _, name := range order {
		clusters = append(clusters, Cluster{Name: name, Resources: groups[name]})
	}
	return clusters
}

// ExtractCluster builds a WADLApplication from the resources in a single Cluster.
func (e *Extractor) ExtractCluster(c Cluster) *WADLApplication {
	app := &WADLApplication{
		Doc:      e.app.Doc,
		Grammars: e.app.Grammars,
		Attrs:    e.app.Attrs,
		Resources: &Resources{
			SampleBase: "http://localhost/sep/",
			Resources:  c.Resources,
		},
	}
	return app
}

// ExtractClusters builds a WADLApplication combining resources from multiple
// clusters. Resources are returned in WADL document order.
func (e *Extractor) ExtractClusters(clusters []Cluster) *WADLApplication {
	// Build a set of resource IDs to include.
	include := make(map[string]bool)
	for _, c := range clusters {
		for _, r := range c.Resources {
			include[r.ID] = true
		}
	}

	resourceMap := extractResourceBlocks(e.rawXML, nil)

	var filtered []Resource
	if e.app.Resources != nil {
		for _, res := range e.app.Resources.Resources {
			if include[res.ID] {
				res.RawXML = resourceMap[res.ID]
				filtered = append(filtered, res)
			}
		}
	}

	return &WADLApplication{
		Doc:      e.app.Doc,
		Grammars: e.app.Grammars,
		Attrs:    e.app.Attrs,
		Resources: &Resources{
			SampleBase: "http://localhost/sep/",
			Resources:  filtered,
		},
	}
}
