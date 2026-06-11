package scaffold

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildSpecWithoutPort(t *testing.T) {
	t.Parallel()

	app := &sep2Application{}
	if err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<application max_entities="42">
  <resources>
    <resource id="BillingReadingSetList" wx:samplePath="/brs" xmlns:wx="urn:ieee:std:2030.5:wadlExt">
      <method id="GETBillingReadingSetList" name="GET" wx:mode="M">
        <response status="200">
          <representation element="sep:BillingReadingSetList"/>
        </response>
      </method>
    </resource>
  </resources>
</application>`), app); err != nil {
		t.Fatalf("unmarshal WADL: %v", err)
	}

	spec, err := buildSpec(app, "BRS")
	if err != nil {
		t.Fatalf("buildSpec returned error without port: %v", err)
	}

	if spec.MaxEntities != 42 {
		t.Fatalf("MaxEntities = %d, want 42", spec.MaxEntities)
	}
	if len(spec.Resources) != 1 {
		t.Fatalf("len(Resources) = %d, want 1", len(spec.Resources))
	}
}

func TestUpdateRoutesAssignsNextPort(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	routesDir := filepath.Join(tmpDir, "internal", "routes")
	if err := os.MkdirAll(routesDir, 0755); err != nil {
		t.Fatalf("mkdir routes dir: %v", err)
	}

	const original = `package routes

import "fmt"

const (
	Core = "egot.internal.com:8000"
	FlowReservation = "egot.internal.com:8001"
	Rsps = "rsps.internal.com:8003"
)

var serviceMap = map[string]string{
	"/dcap": Core,
}

// LinkFor returns the full https URL template for a registered path.
func LinkFor(path string) string {
	host := serviceMap[path]
	return "https://" + host + path
}
`
	routesPath := filepath.Join(routesDir, "routes.go")
	if err := os.WriteFile(routesPath, []byte(original), 0644); err != nil {
		t.Fatalf("write routes.go: %v", err)
	}

	g := &Generator{
		spec: serviceSpec{
			Resources: []resourceSpec{
				{Path: "/new"},
				{Path: "/new/{id1}"},
			},
		},
		serviceName: "new-service",
	}

	if err := g.updateRoutes(tmpDir, "new-service"); err != nil {
		t.Fatalf("updateRoutes: %v", err)
	}

	updated, err := os.ReadFile(routesPath)
	if err != nil {
		t.Fatalf("read routes.go: %v", err)
	}
	content := string(updated)

	if !strings.Contains(content, "NewService") || !strings.Contains(content, `"localhost:8004"`) {
		t.Fatalf("routes.go missing next port assignment:\n%s", content)
	}
	if !strings.Contains(content, `"/new"`) || !strings.Contains(content, `NewService`) {
		t.Fatalf("routes.go missing /new mapping:\n%s", content)
	}
	if !strings.Contains(content, `"/new/{id1}"`) || !strings.Contains(content, `NewService`) {
		t.Fatalf("routes.go missing /new/{id1} mapping:\n%s", content)
	}
}
