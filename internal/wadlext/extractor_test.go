package wadlext

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractClustersPreservesRawXMLForSave(t *testing.T) {
	const wadl = `<?xml version="1.0" encoding="UTF-8"?>
<application xmlns:wx="urn:ieee:std:2030.5:wadlExt" xmlns="http://wadl.dev.java.net/2009/02">
  <resources wx:sampleBase="http://localhost/sep/">
    <resource id="edev-root" wx:samplePath="/edev">
      <method id="getEdevRoot" name="GET"/>
    </resource>
    <resource id="ns-list" wx:samplePath="/edev/{id1}/ns">
      <method id="getNsList" name="GET"/>
    </resource>
    <resource id="der-list" wx:samplePath="/edev/{id1}/der">
      <method id="getDerList" name="GET"/>
    </resource>
  </resources>
</application>
`

	extractor, err := NewExtractorFromReader(strings.NewReader(wadl))
	if err != nil {
		t.Fatalf("NewExtractorFromReader() error = %v", err)
	}

	clusters := extractor.SuggestClusters("/edev", 2)
	selected := make([]Cluster, 0, 2)
	for _, cluster := range clusters {
		if cluster.Name == "ns" || cluster.Name == "der" {
			selected = append(selected, cluster)
		}
	}
	if len(selected) != 2 {
		t.Fatalf("selected %d clusters, want 2", len(selected))
	}

	app := extractor.ExtractClusters(selected)
	if got := app.GetResourceCount(); got != 2 {
		t.Fatalf("ExtractClusters() resource count = %d, want 2", got)
	}

	outputPath := filepath.Join(t.TempDir(), "group.wadl")
	if err := extractor.Save(app, outputPath); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	output := string(data)
	if !strings.Contains(output, `id="ns-list"`) {
		t.Fatalf("output missing ns resource:\n%s", output)
	}
	if !strings.Contains(output, `id="der-list"`) {
		t.Fatalf("output missing der resource:\n%s", output)
	}
	if strings.Contains(output, `id="edev-root"`) {
		t.Fatalf("output unexpectedly included core resource:\n%s", output)
	}
}
