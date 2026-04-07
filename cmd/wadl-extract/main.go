package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Tylores/egot/internal/wadlext"
)

func main() {
	wadlPath := flag.String("wadl", "", "Path to input WADL specification file")
	pathPrefix := flag.String("path", "", "Resource path prefix to extract (e.g., /dr, /edev)")
	outputPath := flag.String("output", "", "Output path for extracted WADL file")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	if *wadlPath == "" || *pathPrefix == "" || *outputPath == "" {
		log.Fatal("Usage: wadl-extract -wadl <input_wadl> -path <path_prefix> -output <output_wadl>")
	}

	// Load the WADL file
	extractor, err := wadlext.NewExtractor(*wadlPath)
	if err != nil {
		log.Fatalf("Failed to load WADL: %v\n", err)
	}

	if *verbose {
		fmt.Printf("✓ Loaded WADL from: %s\n", *wadlPath)
	}

	// Extract resources by path prefix
	extracted := extractor.Extract(*pathPrefix)

	extracted.AddAttribute("port", "8000")
	extracted.AddAttribute("max_entities", "100")

	if extracted.GetResourceCount() == 0 {
		log.Fatalf("No resources found with path prefix: %s\n", *pathPrefix)
	}

	if *verbose {
		fmt.Printf("✓ Found %d resources with prefix '%s'\n", extracted.GetResourceCount(), *pathPrefix)
		for _, path := range extracted.GetResourcePaths() {
			fmt.Printf("  - %s\n", path)
		}
	}

	// Save the extracted WADL
	err = extractor.Save(extracted, *outputPath)
	if err != nil {
		log.Fatalf("Failed to save extracted WADL: %v\n", err)
	}

	fmt.Printf("✓ Extracted WADL saved to: %s\n", *outputPath)
	fmt.Printf("  Resources: %d\n", extracted.GetResourceCount())
	fmt.Printf("  Path prefix: %s\n", *pathPrefix)
}
