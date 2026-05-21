package main

import (
	"encoding/csv"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

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
	startTime := flag.Int64("start", 0, "Start timestamp (unix nanos)")
	endTime := flag.Int64("end", 0, "End timestamp (unix nanos)")
	deviceFilter := flag.String("device", "", "Filter by device SFDI")
	stepSize := flag.Duration("step", 1*time.Minute, "Step size for Step column calculation")
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

	// Header: Timestamp, Step, Device, Power(W)
	writer.Write([]string{"Timestamp", "Step", "Device", "Power(W)"})

	keys := s.Keys()
	sort.Strings(keys)

	var firstTimestamp int64
	count := 0

	for _, key := range keys {
		if !strings.HasPrefix(key, "reading:") {
			continue
		}

		parts := strings.Split(key, ":")
		if len(parts) < 4 {
			continue
		}

		sfdi := parts[1]
		mrid := parts[2]
		ts, _ := strconv.ParseInt(parts[3], 10, 64)

		// Filters
		if *startTime > 0 && ts < *startTime {
			continue
		}
		if *endTime > 0 && ts > *endTime {
			continue
		}
		if *deviceFilter != "" && sfdi != *deviceFilter {
			continue
		}

		if val, ok := s.Get(key); ok {
			if r, ok := val.(*sep.MirrorMeterReading); ok {
				if firstTimestamp == 0 {
					firstTimestamp = ts
				}

				step := (ts - firstTimestamp) / stepSize.Nanoseconds()
				device := sfdi + ":" + mrid
				
				power := "0"
				if r.Reading != nil {
					power = fmt.Sprintf("%d", r.Reading.Value)
				}
				
				writer.Write([]string{
					strconv.FormatInt(ts, 10),
					strconv.FormatInt(step, 10),
					device,
					power,
				})
				count++
			}
		}
	}

	fmt.Printf("✅ Exported %d readings to %s (Start: %v, StepSize: %v)\n", count, *outputPath, time.Unix(0, firstTimestamp), *stepSize)
}
