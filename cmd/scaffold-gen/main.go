package main

import (
	"flag"
	"log"

	"github.com/Tylores/egot/internal/scaffold"
)

func main() {
	wadlPath := flag.String("wadl", "", "Path to WADL specification file")
	outputDir := flag.String("output", "", "Output directory for generated service")
	sep2WadlPath := flag.String("sep2-wadl", "", "Path to SEP 2 WADL file (alternative to -wadl)")
	serviceName := flag.String("service", "", "Service name (required with -sep2-wadl)")
	port := flag.Int("port", 0, "Service port (required with -sep2-wadl)")
	maxEntities := flag.Int("max-entities", 100, "Max entities pool size")
	flag.Parse()

	if *outputDir == "" {
		log.Fatal("Usage: scaffold-gen -wadl <wadl_path> -output <output_dir> OR -sep2-wadl <sep2_wadl_path> -service <name> -port <port> -output <output_dir>")
	}

	var gen *scaffold.Generator
	var err error

	// Use SEP 2 WADL if provided, otherwise use regular WADL
	if *sep2WadlPath != "" {
		if *serviceName == "" || *port == 0 {
			log.Fatal("When using -sep2-wadl, both -service and -port are required")
		}
		gen, err = scaffold.NewGeneratorFromSEP2WADL(*sep2WadlPath, *serviceName, *port, *maxEntities)
	} else if *wadlPath != "" {
		gen, err = scaffold.NewGenerator(*wadlPath)
	} else {
		log.Fatal("Either -wadl or -sep2-wadl must be provided")
	}

	if err != nil {
		log.Fatalf("Failed to load WADL: %v", err)
	}

	err = gen.Generate(*outputDir)
	if err != nil {
		log.Fatalf("Failed to generate scaffold: %v", err)
	}

	log.Printf("✓ Scaffold generated successfully in %s\n", *outputDir)
	log.Printf("Next steps:\n")
	log.Printf("  1. Add business logic to internal/%s/handler/handler.go\n", gen.ServiceName())
	log.Printf("  2. Implement storage in internal/%s/repository/memory/memory.go\n", gen.ServiceName())
	log.Printf("  3. Run: go build ./cmd/%s\n", gen.ServiceName())
}
