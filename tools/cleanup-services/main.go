package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PRODUCTION_SERVICES lists all official/production microservices to preserve
var PRODUCTION_SERVICES = map[string]bool{
	"BRS":             false,
	"Bill":            false,
	"DCAP":            false,
	"DERP":            false,
	"DR":              false,
	"EDevice":         false,
	"File":            false,
	"MUP":             false,
	"Messaging":       false,
	"Notify":          false,
	"PPY":             false,
	"SDevice":         false,
	"TariffProfile":   false,
	"TimeOfUse":       false,
	"UPT":             false,
	"DER":             false,
	"FlowReservation": false,
	"core":            false,
	"der":             false,
	"rsps":            false,
	"flowreservation": false,
}

type ServiceInfo struct {
	Name    string
	CmdPath string
	IntPath string
	IsStale bool
	Size    int64
}

func main() {
	var (
		dryRun       = flag.Bool("dry-run", true, "Show what would be deleted without actually deleting")
		interactive  = flag.Bool("interactive", true, "Require confirmation before deletion")
		removeAll    = flag.Bool("remove-all-stale", false, "Remove all stale services without prompting")
		listServices = flag.Bool("list", false, "List all services and their status")
	)
	flag.Parse()

	repoRoot := getRepoRoot()
	services := scanServices(repoRoot)

	if *listServices {
		displayServicesList(services)
		return
	}

	staleServices := filterStaleServices(services)
	if len(staleServices) == 0 {
		fmt.Println("✓ No stale services found!")
		return
	}

	fmt.Printf("\nFound %d stale service(s):\n\n", len(staleServices))
	for _, svc := range staleServices {
		fmt.Printf("  • %s\n", svc.Name)
	}

	if !*interactive && !*removeAll {
		fmt.Println("\nUse --interactive=false --remove-all-stale to proceed without confirmation")
		return
	}

	if *interactive && !*removeAll {
		if !confirmRemoval(staleServices) {
			fmt.Println("Cancelled.")
			return
		}
	}

	if *dryRun {
		fmt.Println("\n[DRY RUN] Would remove:")
		for _, svc := range staleServices {
			fmt.Printf("  - %s/\n", svc.CmdPath)
			fmt.Printf("  - %s/\n", svc.IntPath)
		}
		fmt.Println("\nRun with --dry-run=false to actually delete")
		fmt.Println("\nAlso requires manual update to internal/routes/routes.go to remove:")
		for _, svc := range staleServices {
			fmt.Printf("  - %s constant and route entries\n", svc.Name)
		}
		return
	}

	// Actual deletion
	fmt.Println("\nRemoving stale services...")
	deleteServices(repoRoot, staleServices)
	fmt.Println("\n✓ Cleanup complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Manually edit internal/routes/routes.go to remove service constants and serviceMap entries")
	fmt.Println("  2. Run: go build ./...")
	fmt.Println("  3. Commit changes: git add -A && git commit -m \"cleanup: remove stale microservices\"")
}

func getRepoRoot() string {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Look for go.mod to confirm we're in repo root
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
		return root
	}
	fmt.Fprintf(os.Stderr, "Error: Not in egot repository root (no go.mod found)\n")
	os.Exit(1)
	return ""
}

func scanServices(repoRoot string) []ServiceInfo {
	var services []ServiceInfo

	// Scan cmd/
	cmdDir := filepath.Join(repoRoot, "cmd")
	entries, _ := os.ReadDir(cmdDir)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		_, found := PRODUCTION_SERVICES[name]
		if !found {
			continue
		}
		cmdPath := filepath.Join("cmd", name)
		intPath := filepath.Join("internal", name)

		// Calculate size (rough estimate)
		var size int64
		filepath.Walk(filepath.Join(repoRoot, intPath), func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				size += info.Size()
			}
			return nil
		})

		services = append(services, ServiceInfo{
			Name:    name,
			CmdPath: cmdPath,
			IntPath: intPath,
			IsStale: !PRODUCTION_SERVICES[name],
			Size:    size,
		})
	}

	return services
}

func filterStaleServices(services []ServiceInfo) []ServiceInfo {
	var stale []ServiceInfo
	for _, svc := range services {
		if svc.IsStale {
			stale = append(stale, svc)
		}
	}
	return stale
}

func displayServicesList(services []ServiceInfo) {
	fmt.Println("Production Services (PROTECTED):")
	var prodCount int
	for _, svc := range services {
		if !svc.IsStale {
			fmt.Printf("  ✓ %15s  %s  (%d bytes)\n", svc.Name, svc.IntPath, svc.Size)
			prodCount++
		}
	}

	fmt.Println("\nStale Services (Can be removed):")
	var staleCount int
	for _, svc := range services {
		if svc.IsStale {
			fmt.Printf("  ✗ %15s  %s  (%d bytes)\n", svc.Name, svc.IntPath, svc.Size)
			staleCount++
		}
	}

	fmt.Printf("\nTotal: %d production + %d stale = %d services\n", prodCount, staleCount, len(services))
}

func confirmRemoval(services []ServiceInfo) bool {
	fmt.Print("\nProceed with removal? (yes/no): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.ToLower(strings.TrimSpace(scanner.Text())) == "yes"
	}
	return false
}

func deleteServices(repoRoot string, services []ServiceInfo) {
	for _, svc := range services {
		// Remove cmd/{Service}
		cmdPath := filepath.Join(repoRoot, svc.CmdPath)
		if err := os.RemoveAll(cmdPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing %s: %v\n", svc.CmdPath, err)
		} else {
			fmt.Printf("✓ Removed %s/\n", svc.CmdPath)
		}

		// Remove internal/{Service}
		intPath := filepath.Join(repoRoot, svc.IntPath)
		if err := os.RemoveAll(intPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing %s: %v\n", svc.IntPath, err)
		} else {
			fmt.Printf("✓ Removed %s/\n", svc.IntPath)
		}
	}

	// Generate routes.go removal guide
	routesPath := filepath.Join(repoRoot, "internal/routes/routes.go")
	fmt.Printf("\n📝 Manual cleanup required in routes.go:\n")
	fmt.Printf("   Location: %s\n\n", routesPath)
	fmt.Println("   Remove these constant definitions:")
	for _, svc := range services {
		fmt.Printf("     - %s = \"...\"  (look for this constant in const block)\n", svc.Name)
	}
	fmt.Println("\n   Remove these serviceMap entries:")
	for _, svc := range services {
		fmt.Printf("     - Entries with route to %s\n", svc.Name)
	}
}

func isCapitalized(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}
