package wadlext

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// SplitConfig is the top-level structure of a wadl-extract grouping config file.
type SplitConfig struct {
	WADL   string                 `yaml:"wadl"`
	Base   string                 `yaml:"base"`
	Depth  int                    `yaml:"depth"`
	Output string                 `yaml:"output"`
	Groups map[string]GroupConfig `yaml:"groups"`
}

// GroupConfig describes one logical service group within a SplitConfig.
type GroupConfig struct {
	Clusters []string `yaml:"clusters"`
}

// LoadSplitConfig reads and unmarshals a YAML config file.
func LoadSplitConfig(path string) (*SplitConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg SplitConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Depth == 0 {
		cfg.Depth = 2
	}

	return &cfg, nil
}

// ValidateConfig checks that:
//   - every cluster key referenced in the config exists in allClusters
//   - no cluster key appears in more than one group
//   - every cluster in allClusters is assigned to exactly one group
//
// Returns a non-nil error listing all problems found.
func ValidateConfig(cfg *SplitConfig, allClusters []Cluster) error {
	// Build set of known cluster keys.
	known := make(map[string]bool, len(allClusters))
	for _, c := range allClusters {
		known[c.Name] = true
	}

	// Track which clusters have been assigned and to which group.
	assigned := make(map[string]string) // clusterKey → groupName
	var errs []string

	for groupName, gc := range cfg.Groups {
		for _, key := range gc.Clusters {
			if !known[key] {
				available := clusterNames(allClusters)
				errs = append(errs, fmt.Sprintf(
					"cluster key %q in group %q not found (available: %s)",
					key, groupName, strings.Join(available, ", "),
				))
				continue
			}
			if prev, dup := assigned[key]; dup {
				errs = append(errs, fmt.Sprintf(
					"cluster key %q assigned to both %q and %q",
					key, prev, groupName,
				))
				continue
			}
			assigned[key] = groupName
		}
	}

	// Find clusters that weren't assigned.
	var unassigned []string
	for _, c := range allClusters {
		if _, ok := assigned[c.Name]; !ok {
			unassigned = append(unassigned, c.Name)
		}
	}
	if len(unassigned) > 0 {
		sort.Strings(unassigned)
		errs = append(errs, fmt.Sprintf(
			"unassigned clusters: %s — every cluster must be in a group",
			strings.Join(unassigned, ", "),
		))
	}

	if len(errs) > 0 {
		sort.Strings(errs)
		return fmt.Errorf("config validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

// GenerateInitConfig returns YAML text for a starter config file with all
// clusters listed under an "# unassigned" stub, ready for the user to move
// into named groups.
func GenerateInitConfig(clusters []Cluster, wadlPath, base, output string, depth int) string {
	var sb strings.Builder

	sb.WriteString("wadl: ")
	sb.WriteString(wadlPath)
	sb.WriteString("\nbase: ")
	sb.WriteString(base)
	sb.WriteString("\ndepth: ")
	sb.WriteString(fmt.Sprintf("%d", depth))
	sb.WriteString("\noutput: ")
	if output == "" {
		output = "wadl/output/"
	}
	sb.WriteString(output)
	sb.WriteString("\n\n")

	sb.WriteString("# Define your groups below. Each group becomes one WADL file (<name>.wadl)\n")
	sb.WriteString("# and one microservice. Every cluster key must be assigned to exactly one group.\n")
	sb.WriteString("#\n")
	sb.WriteString("# Available cluster keys:\n")
	for _, c := range clusters {
		sb.WriteString(fmt.Sprintf("#   %-14s  (%d paths)\n", c.Name, len(c.Resources)))
	}
	sb.WriteString("\ngroups:\n")
	sb.WriteString("  # Example group — rename and adjust clusters as needed:\n")
	sb.WriteString("  # my-service:\n")
	sb.WriteString("  #   clusters:\n")
	for _, c := range clusters {
		sb.WriteString(fmt.Sprintf("  #     - %s\n", c.Name))
	}

	return sb.String()
}

// clusterNames returns a sorted slice of cluster names from a Cluster slice.
func clusterNames(clusters []Cluster) []string {
	names := make([]string, len(clusters))
	for i, c := range clusters {
		names[i] = c.Name
	}
	sort.Strings(names)
	return names
}
