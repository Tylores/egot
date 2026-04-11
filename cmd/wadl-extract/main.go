package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Tylores/egot/internal/wadlext"
)

func main() {
	wadlPath := flag.String("wadl", "", "Path to input WADL specification file")
	pathArg := flag.String("path", "", "Comma-separated resource path prefixes to extract (e.g., /dr or /edev/{id1}/der,/edev/{id1}/ns)")
	outputPath := flag.String("output", "", "Output path for extracted WADL file")
	suggest := flag.Bool("suggest", false, "Print cluster suggestions for -path at -depth (no files written)")
	splitDir := flag.String("split", "", "Directory to write one WADL per cluster (uses -path and -depth)")
	depth := flag.Int("depth", 2, "Clustering depth for -suggest, -split, and -init-config (default: 2)")
	configPath := flag.String("config", "", "Path to a YAML grouping config file; generates one WADL per group")
	initConfig := flag.String("init-config", "", "Write a starter grouping config to this path (requires -wadl and -path)")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	// Config-file mode: only needs -config
	if *configPath != "" {
		runConfig(*configPath, *verbose)
		return
	}

	// Init-config mode: needs -wadl and -path
	if *initConfig != "" {
		if *wadlPath == "" || *pathArg == "" {
			log.Fatal("-init-config requires -wadl and -path")
		}
		extractor, err := wadlext.NewExtractor(*wadlPath)
		if err != nil {
			log.Fatalf("Failed to load WADL: %v\n", err)
		}
		prefixes := splitPaths(*pathArg)
		if len(prefixes) != 1 {
			log.Fatal("-init-config requires a single -path prefix")
		}
		runInitConfig(extractor, *wadlPath, prefixes[0], *depth, *initConfig)
		return
	}

	// Standard extract / suggest / split modes all need -wadl and -path
	if *wadlPath == "" || *pathArg == "" {
		log.Fatal("Usage: wadl-extract -wadl <input_wadl> -path <prefix>[,<prefix>...] [-output <output_wadl>] [-suggest] [-split <dir>] [-depth <n>]\n" +
			"       wadl-extract -wadl <input_wadl> -path <prefix> -init-config <config.yaml> [-depth <n>]\n" +
			"       wadl-extract -config <config.yaml>")
	}

	if !*suggest && *splitDir == "" && *outputPath == "" {
		log.Fatal("-output is required unless -suggest or -split is used")
	}

	prefixes := splitPaths(*pathArg)

	extractor, err := wadlext.NewExtractor(*wadlPath)
	if err != nil {
		log.Fatalf("Failed to load WADL: %v\n", err)
	}

	if *verbose {
		fmt.Printf("✓ Loaded WADL from: %s\n", *wadlPath)
	}

	switch {
	case *suggest:
		if len(prefixes) > 1 {
			log.Fatal("-suggest only supports a single -path prefix")
		}
		runSuggest(extractor, prefixes[0], *depth)
	case *splitDir != "":
		if len(prefixes) > 1 {
			log.Fatal("-split only supports a single -path prefix")
		}
		runSplit(extractor, prefixes[0], *depth, *splitDir, *verbose)
	default:
		runExtract(extractor, prefixes, *outputPath, *verbose)
	}
}

// splitPaths parses a comma-separated path list, trimming whitespace.
func splitPaths(arg string) []string {
	parts := strings.Split(arg, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func runExtract(extractor *wadlext.Extractor, prefixes []string, outputPath string, verbose bool) {
	extracted := extractor.ExtractMany(prefixes)

	if extracted.GetResourceCount() == 0 {
		log.Fatalf("No resources found for paths: %s\n", strings.Join(prefixes, ", "))
	}

	if verbose {
		fmt.Printf("✓ Found %d resources for %d prefix(es)\n", extracted.GetResourceCount(), len(prefixes))
		for _, path := range extracted.GetResourcePaths() {
			fmt.Printf("  - %s\n", path)
		}
	}

	if err := extractor.Save(extracted, outputPath); err != nil {
		log.Fatalf("Failed to save extracted WADL: %v\n", err)
	}

	fmt.Printf("✓ Extracted WADL saved to: %s\n", outputPath)
	fmt.Printf("  Resources: %d\n", extracted.GetResourceCount())
	fmt.Printf("  Path prefix(es): %s\n", strings.Join(prefixes, ", "))
}

func runSuggest(extractor *wadlext.Extractor, pathPrefix string, depth int) {
	clusters := extractor.SuggestClusters(pathPrefix, depth)
	if len(clusters) == 0 {
		fmt.Printf("No resources found with path prefix: %s\n", pathPrefix)
		return
	}

	total := 0
	for _, c := range clusters {
		total += len(c.Resources)
	}

	fmt.Printf("Clusters for %s at depth %d (%d clusters, %d total paths):\n", pathPrefix, depth, len(clusters), total)
	for _, c := range clusters {
		paths := c.Paths()
		preview := ""
		if len(paths) > 0 {
			preview = paths[0]
			if len(paths) > 1 {
				preview += ", ..."
			}
		}
		fmt.Printf("  [%-12s] %3d paths  → %s\n", c.Name, len(c.Resources), preview)
	}
}

func runSplit(extractor *wadlext.Extractor, pathPrefix string, depth int, outDir string, verbose bool) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v\n", err)
	}

	clusters := extractor.SuggestClusters(pathPrefix, depth)
	if len(clusters) == 0 {
		log.Fatalf("No resources found with path prefix: %s\n", pathPrefix)
	}

	fmt.Printf("Splitting %s (depth %d) into %d clusters → %s\n", pathPrefix, depth, len(clusters), outDir)
	for _, c := range clusters {
		app := extractor.ExtractCluster(c)

		filename := filepath.Join(outDir, c.Name+".wadl")
		if err := extractor.Save(app, filename); err != nil {
			log.Fatalf("Failed to write cluster %s: %v\n", c.Name, err)
		}

		if verbose {
			for _, p := range c.Paths() {
				fmt.Printf("    %s\n", p)
			}
		}
		fmt.Printf("  ✓ [%s] %d paths → %s\n", c.Name, len(c.Resources), filename)
	}
}

func runConfig(configPath string, verbose bool) {
	cfg, err := wadlext.LoadSplitConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	extractor, err := wadlext.NewExtractor(cfg.WADL)
	if err != nil {
		log.Fatalf("Failed to load WADL: %v\n", err)
	}

	clusters := extractor.SuggestClusters(cfg.Base, cfg.Depth)
	if len(clusters) == 0 {
		log.Fatalf("No clusters found for base %q at depth %d\n", cfg.Base, cfg.Depth)
	}

	// Validate before writing anything.
	if err := wadlext.ValidateConfig(cfg, clusters); err != nil {
		log.Fatalf("%v\n", err)
	}

	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v\n", err)
	}

	// Build a lookup: cluster key → Cluster
	clusterMap := make(map[string]wadlext.Cluster, len(clusters))
	for _, c := range clusters {
		clusterMap[c.Name] = c
	}

	// Sort group names for deterministic output.
	groupNames := make([]string, 0, len(cfg.Groups))
	for name := range cfg.Groups {
		groupNames = append(groupNames, name)
	}
	sortStrings(groupNames)

	fmt.Printf("Generating %d group WADLs from %s → %s\n", len(cfg.Groups), cfg.WADL, cfg.Output)

	for _, groupName := range groupNames {
		gc := cfg.Groups[groupName]

		// Collect the Cluster objects for this group.
		groupClusters := make([]wadlext.Cluster, 0, len(gc.Clusters))
		for _, key := range gc.Clusters {
			groupClusters = append(groupClusters, clusterMap[key])
		}

		app := extractor.ExtractClusters(groupClusters)

		filename := filepath.Join(cfg.Output, groupName+".wadl")
		if err := extractor.Save(app, filename); err != nil {
			log.Fatalf("Failed to write group %s: %v\n", groupName, err)
		}

		if verbose {
			for _, p := range app.GetResourcePaths() {
				fmt.Printf("    %s\n", p)
			}
		}
		fmt.Printf("  ✓ [%-20s] %2d paths → %s\n",
			groupName, app.GetResourceCount(), filename)
	}
}

func runInitConfig(extractor *wadlext.Extractor, wadlPath, base string, depth int, outputPath string) {
	clusters := extractor.SuggestClusters(base, depth)
	if len(clusters) == 0 {
		log.Fatalf("No clusters found for base %q at depth %d\n", base, depth)
	}

	content := wadlext.GenerateInitConfig(clusters, wadlPath, base, "", depth)
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write init config: %v\n", err)
	}

	fmt.Printf("✓ Starter config written to: %s\n", outputPath)
	fmt.Printf("  %d clusters listed (%s at depth %d)\n", len(clusters), base, depth)
	fmt.Printf("  Edit the file to define groups, then run:\n")
	fmt.Printf("    wadl-extract -config %s\n", outputPath)
}

// sortStrings sorts a string slice in place.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
