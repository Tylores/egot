package scaffold

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// SEP2WADLApplication represents a SEP 2 WADL application root
type SEP2WADLApplication struct {
	XMLName   xml.Name           `xml:"application"`
	Resources *SEP2Resources     `xml:"resources"`
	Attrs     []xml.Attr         `xml:",any,attr"`
}

// SEP2Resources represents the resources section of SEP 2 WADL
type SEP2Resources struct {
	Resources []SEP2Resource `xml:"resource"`
	Attrs     []xml.Attr     `xml:",any,attr"`
}

// SEP2Resource represents a resource in SEP 2 WADL
type SEP2Resource struct {
	XMLName     xml.Name       `xml:"resource"`
	ID          string         `xml:"id,attr"`
	SamplePath  string         `xml:""`
	Doc         *SEP2Doc       `xml:"doc"`
	Methods     []SEP2Method   `xml:"method"`
	SampleParam []SEP2SampleParam `xml:""`
	Attrs       []xml.Attr     `xml:",any,attr"`
}

// UnmarshalXML handles namespaced attributes in SEP 2 WADL
func (r *SEP2Resource) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias SEP2Resource
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	// Extract samplePath and sampleParam from namespaced attributes
	for _, attr := range start.Attr {
		if attr.Name.Space == "urn:ieee:std:2030.5:wadlExt" && attr.Name.Local == "samplePath" {
			r.SamplePath = attr.Value
		}
	}

	return d.DecodeElement(aux, &start)
}

// SEP2Doc represents documentation in SEP 2 WADL
type SEP2Doc struct {
	Title string `xml:"title,attr"`
	Text  string `xml:",chardata"`
}

// SEP2Method represents a method in SEP 2 WADL
type SEP2Method struct {
	XMLName  xml.Name        `xml:"method"`
	ID       string          `xml:"id,attr"`
	Name     string          `xml:"name,attr"`
	Request  *SEP2Request    `xml:"request"`
	Response []SEP2Response  `xml:"response"`
	Attrs    []xml.Attr      `xml:",any,attr"`
}

// SEP2Request represents a request in SEP 2 WADL
type SEP2Request struct {
	Params []SEP2Param           `xml:"param"`
	Reps   []SEP2Representation  `xml:"representation"`
}

// SEP2Response represents a response in SEP 2 WADL
type SEP2Response struct {
	Status string                `xml:"status,attr"`
	Params []SEP2Param           `xml:"param"`
	Reps   []SEP2Representation  `xml:"representation"`
}

// SEP2Param represents a parameter in SEP 2 WADL
type SEP2Param struct {
	Name     string `xml:"name,attr"`
	Style    string `xml:"style,attr"`
	Type     string `xml:"type,attr"`
	Required string `xml:"required,attr"`
}

// SEP2Representation represents a media type representation in SEP 2 WADL
type SEP2Representation struct {
	MediaType string `xml:"mediaType,attr"`
	Element   string `xml:"element,attr"`
}

// SEP2SampleParam represents a sample parameter in SEP 2 WADL
type SEP2SampleParam struct {
	XMLName xml.Name   `xml:""`
	Name    string     `xml:"name,attr"`
	Style   string     `xml:"style,attr"`
	Type    string     `xml:"type,attr"`
}

// UnmarshalXML handles sampleParam elements
func (p *SEP2SampleParam) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local == "sampleParam" && start.Name.Space == "urn:ieee:std:2030.5:wadlExt" {
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "name":
				p.Name = attr.Value
			case "style":
				p.Style = attr.Value
			case "type":
				p.Type = attr.Value
			}
		}
	}
	d.Skip()
	return nil
}

// ConvertSEP2WADL converts a SEP 2 WADL specification to the scaffold generator format
func ConvertSEP2WADL(sep2WadlPath, serviceName string, port, maxEntities int) (*WADL, error) {
	data, err := os.ReadFile(sep2WadlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SEP 2 WADL file: %w", err)
	}

	sep2App := &SEP2WADLApplication{}
	err = xml.Unmarshal(data, sep2App)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SEP 2 WADL: %w", err)
	}

	// Create scaffold WADL
	wadl := &WADL{
		Name:        serviceName,
		Port:        port,
		MaxEntities: maxEntities,
		Resources:   []Resource{},
	}

	if sep2App.Resources == nil || len(sep2App.Resources.Resources) == 0 {
		return nil, fmt.Errorf("SEP 2 WADL has no resources")
	}

	// Convert resources
	for _, sep2Res := range sep2App.Resources.Resources {
		if sep2Res.SamplePath == "" {
			continue
		}

		scaffoldRes := Resource{
			Path:    sep2Res.SamplePath,
			Methods: []Method{},
		}

		// Extract path parameters from sample path (e.g., {id1}, {id2})
		pathParamNames := extractPathParams(sep2Res.SamplePath)
		pathParamMap := make(map[string]string)
		for _, name := range pathParamNames {
			pathParamMap[name] = "string" // Default to string type
		}

		// Convert methods
		for _, sep2Method := range sep2Res.Methods {
			httpMethod := sep2Method.Name
			methodName := sep2Method.ID

			// Extract response type from response representation
			responseType := ""
			for _, resp := range sep2Method.Response {
				for _, rep := range resp.Reps {
					if rep.Element != "" {
						// Extract local name from element (e.g., "sep:ResponseSet" -> "ResponseSet")
						parts := strings.Split(rep.Element, ":")
						if len(parts) > 0 {
							responseType = parts[len(parts)-1]
						}
						break
					}
				}
				if responseType != "" {
					break
				}
			}

			// Extract request type from request representation
			requestType := ""
			if sep2Method.Request != nil {
				for _, rep := range sep2Method.Request.Reps {
					if rep.Element != "" {
						// Extract local name from element
						parts := strings.Split(rep.Element, ":")
						if len(parts) > 0 {
							requestType = parts[len(parts)-1]
						}
						break
					}
				}
			}

			scaffoldMethod := Method{
				Name:         methodName,
				HTTPMethod:   httpMethod,
				ResponseType: responseType,
				RequestType:  requestType,
				PathParams:   []PathParam{},
			}

			// Add path parameters
			for _, paramName := range pathParamNames {
				scaffoldMethod.PathParams = append(scaffoldMethod.PathParams, PathParam{
					Name: paramName,
					Type: pathParamMap[paramName],
				})
			}

			scaffoldRes.Methods = append(scaffoldRes.Methods, scaffoldMethod)
		}

		if len(scaffoldRes.Methods) > 0 {
			wadl.Resources = append(wadl.Resources, scaffoldRes)
		}
	}

	return wadl, nil
}

// extractPathParams extracts parameter names from a path string
// e.g., "/rsps/{id1}/rsp/{id2}" -> ["id1", "id2"]
func extractPathParams(path string) []string {
	var params []string
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName := strings.TrimPrefix(strings.TrimSuffix(part, "}"), "{")
			params = append(params, paramName)
		}
	}
	return params
}
