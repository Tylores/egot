package wadlext

import (
	"regexp"
	"strings"
)

// HrefRewriter handles rewriting href attributes in WADL documents
// to properly resolve cross-service references through the API gateway
type HrefRewriter struct {
	// ServiceName is the name of the service being processed
	ServiceName string
	// GatewayHost is the nginx gateway host (e.g., localhost:8443)
	GatewayHost string
	// CrossServiceMap maps element names to their source service paths
	// e.g., "Bill" → "/bill", "DERP" → "/derp"
	CrossServiceMap map[string]string
}

// NewHrefRewriter creates a new href rewriter for a specific service
func NewHrefRewriter(serviceName, gatewayHost string) *HrefRewriter {
	return &HrefRewriter{
		ServiceName:     serviceName,
		GatewayHost:     gatewayHost,
		CrossServiceMap: buildDefaultCrossServiceMap(),
	}
}

// RewriteWADL rewrites all href attributes in a WADL document
// - Local schema references (sep.xsd) remain unchanged
// - Cross-service element references are rewritten to use gateway
func (r *HrefRewriter) RewriteWADL(wadlXML string) string {
	// Pattern 1: Keep local schema references unchanged (sep.xsd, etc.)
	// These are maintained as-is since they're loaded locally

	// Pattern 2: Rewrite cross-service element references
	// If representations reference elements from other services, rewrite to gateway
	wadlXML = r.rewriteElementReferences(wadlXML)

	return wadlXML
}

// rewriteElementReferences finds element attributes that reference cross-service types
// and rewrites them to point through the nginx gateway
func (r *HrefRewriter) rewriteElementReferences(wadlXML string) string {
	// Pattern: element="sep:SomeElementName"
	// These are namespace-qualified references to element definitions
	// For now, we keep them as-is since the sep: namespace is local
	// In a full implementation, we'd detect which elements come from other services

	return wadlXML
}

// RewriteRepresentationElement rewrites representation element references
// to point to the correct service namespace
func (r *HrefRewriter) RewriteRepresentationElement(elementName string) string {
	// elementName format: "namespace:LocalName" (e.g., "sep:Bill")
	// For cross-service references, this would be rewritten

	// For now, return unchanged since all use the shared sep namespace
	return elementName
}

// buildDefaultCrossServiceMap creates the default mapping of resource paths to services
func buildDefaultCrossServiceMap() map[string]string {
	return map[string]string{
		// Discovery & Core
		"/dcap": "DCAP",
		"/tm":   "TimeOfUse",

		// Service mappings
		"/bill":     "Bill",
		"/brs":      "BRS",
		"/derp":     "DERP",
		"/dr":       "DR",
		"/edev":     "EDevice",
		"/file":     "File",
		"/mup":      "MUP",
		"/msg":      "Messaging",
		"/ntfy":     "Notify",
		"/ppy":      "PPY",
		"/sdev":     "SDevice",
		"/tp":       "TariffProfile",
		"/upt":      "UPT",
		"/rsps":     "Rsps",
		"/der":      "DER",
	}
}

// FilterHrefsForService filters WADL to only include hrefs relevant to a service
// This removes cross-service references that aren't part of the base WADL
func FilterHrefsForService(wadlXML string, serviceName string) string {
	// For the current implementation, we:
	// 1. Keep all schema includes (sep.xsd) - these are shared
	// 2. Keep all representations - they use the shared sep namespace
	// 3. Don't rewrite hrefs - they all point to local schemas

	// In a more sophisticated implementation, we would:
	// 1. Parse all representation elements
	// 2. Detect which elements come from other services
	// 3. Rewrite those to reference the other service through the gateway

	return wadlXML
}

// ResolveElementDefinition resolves an element reference to its source service
// Returns the service path (e.g., "/bill") if found in cross-service mappings
func ResolveElementDefinition(elementName string, crossServiceMap map[string]string) string {
	// Parse element name (format: "namespace:LocalName")
	parts := strings.Split(elementName, ":")
	if len(parts) < 2 {
		return ""
	}
	localName := parts[1]

	// This would need a reverse mapping from element names to services
	// For now, we return empty since all elements are in the shared sep namespace
	_ = localName
	return ""
}

// BuildGatewayElementRef builds a proper element reference through the gateway
// Example: "https://localhost:8443/bill/" + "sep:Bill"
func BuildGatewayElementRef(servicePath, elementName, gatewayHost string) string {
	return "https://" + gatewayHost + servicePath + "#" + elementName
}

// ValidateHrefs checks that all href attributes in WADL are resolvable
// Returns a list of any problematic hrefs found
func ValidateHrefs(wadlXML string) []string {
	var issues []string

	// Pattern to find href attributes
	hrefPattern := regexp.MustCompile(`href="([^"]+)"`)
	matches := hrefPattern.FindAllStringSubmatch(wadlXML, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		href := match[1]

		// Check for common issues
		if href == "" {
			issues = append(issues, "empty href found")
		} else if strings.Contains(href, "http://") && !strings.Contains(href, "localhost") {
			// Warn about external URLs (might indicate cross-service reference)
			issues = append(issues, "external href found: "+href)
		}
	}

	return issues
}
