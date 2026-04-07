package main

import (
	"flag"
	"log"

	"github.com/Tylores/egot/internal/scaffold"
)

func main() {
	wadlPath := flag.String("wadl", "", "Path to WADL specification file")
	flag.Parse()

	if *wadlPath == "" {
		log.Fatal("Usage: scaffold-gen -wadl <wadl_path>")
	}

	gen, err := scaffold.NewGenerator(*wadlPath)
	if err != nil {
		log.Fatalf("Failed to load WADL: %v", err)
	}

	err = gen.Generate()
	if err != nil {
		log.Fatalf("Failed to generate scaffold: %v", err)
	}

	log.Printf("✓ Scaffold generated successfully\n")
	log.Printf("Next steps:\n")
	log.Printf("  1. Add business logic to internal/%s/handler/handler.go\n", gen.ServiceName())
	log.Printf("  2. Implement storage in internal/%s/repository/memory/memory.go\n", gen.ServiceName())
	log.Printf("  3. Run: go build ./cmd/%s\n", gen.ServiceName())
}
