package main

import (
	"encoding/csv"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/sep"
)

func init() {
	gob.Register(&sep.MirrorUsagePoint{})
	gob.Register(&sep.MirrorMeterReading{})
}

func main() {
	dbPath := flag.String("db", "data/mup.db", "Path to mup.db")
	outputPath := flag.String("output", "der_results.csv", "Output CSV path")
	flag.Parse()

	s := store.New(*dbPath)
	if err := s.Load(); err != nil {
		log.Fatalf("Failed to load store: %v", err)
	}

	f, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Header: Timestamp, Device, Power(W)
	writer.Write([]string{"Timestamp", "Device", "Power(W)"})

	keys := s.Keys()
	sort.Strings(keys)

	for _, key := range keys {
		if strings.HasPrefix(key, "reading:") {
			if val, ok := s.Get(key); ok {
				if r, ok := val.(*sep.MirrorMeterReading); ok {
					// key format: reading:sfdi:mrid:nanos
					parts := strings.Split(key, ":")
					if len(parts) < 4 {
						continue
					}
					device := parts[1] + ":" + parts[2]
					timestamp := parts[3]
					
					power := "0"
					if r.Reading != nil {
						power = fmt.Sprintf("%d", r.Reading.Value)
					}
					
					writer.Write([]string{timestamp, device, power})
				}
			}
		}
	}

	fmt.Printf("✅ Exported %d readings to %s\n", len(keys), *outputPath)
}
