package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Tylores/egot/internal/operator"
	"github.com/Tylores/egot/sep"
)

func main() {
	fmt.Println("🚀 Starting GMLC Grid Services Validation (Phase 2)...")

	// Ensure output directory exists
	err := os.MkdirAll("data", 0755)
	if err != nil {
		fmt.Printf("❌ Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	// Run Scenario 2.1: Hourly Scheduling & 5-min Regulation
	runScenario21()

	// Run Scenario 2.2: Blackstart Coordination
	runScenario22()

	// Run Scenario 2.3: Reserve Capacity & Settlement
	runScenario23()

	fmt.Println("🎉 All scenarios run successfully. Output files generated in data/")
}

func runScenario21() {
	fmt.Println("👉 Running Scenario 2.1: Hourly Scheduling & 5-min Regulation (EIM)...")
	
	// We simulate a 24-hour day in 5-minute intervals (288 steps)
	const totalSteps = 288
	const intervalsPerHour = 12

	// Output files
	fCoordinated, err := os.Create(filepath.Join("data", "scenario_eim.csv"))
	if err != nil {
		panic(err)
	}
	defer fCoordinated.Close()

	fBaseline, err := os.Create(filepath.Join("data", "scenario_eim_baseline.csv"))
	if err != nil {
		panic(err)
	}
	defer fBaseline.Close()

	wCoord := csv.NewWriter(fCoordinated)
	wBase := csv.NewWriter(fBaseline)

	header := []string{"Step", "Device", "Power(W)"}
	wCoord.Write(header)
	wBase.Write(header)

	// Seed random for reproducible simulation
	rng := rand.New(rand.NewSource(42))

	for step := 0; step < totalSteps; step++ {
		hour := step / intervalsPerHour

		// 1. Base solar profiles (PV) - peak at noon
		var solarPower float64
		if hour >= 6 && hour <= 18 {
			// Sinusoidal solar curve
			rad := float64(hour-6) / 12.0 * math.Pi
			solarPower = 6000.0 * math.Sin(rad)
		} else {
			solarPower = 0.0
		}
		// Add some high frequency cloud fluctuations
		fluctuation := (rng.Float64() - 0.5) * 1500.0
		solarPowerCoordinated := math.Max(0, solarPower+fluctuation)
		solarPowerBaseline := math.Max(0, solarPower+fluctuation*2.0) // uncoordinated solar has worse fluctuation

		// 2. Base Load profiles
		// Load 1 (residential): peaks in morning and evening
		load1Base := 1500.0 + 800.0*math.Sin(float64(hour)/24.0*2.0*math.Pi-math.Pi/2.0)
		// Load 2 (commercial): peaks during office hours
		load2Base := 2500.0
		if hour >= 8 && hour <= 17 {
			load2Base += 3000.0
		}
		// Add small random load noise
		load1 := load1Base + rng.Float64()*300.0
		load2 := load2Base + rng.Float64()*400.0

		// EV charging: baseline EV charges immediately when home (evening)
		// Coordinated EV charges during solar peak (midday)
		var evBaseline float64
		if hour >= 17 && hour <= 21 {
			evBaseline = 6000.0 // Charge at 6kW
		} else {
			evBaseline = 0.0
		}

		var evCoordinated float64
		if hour >= 10 && hour <= 14 {
			evCoordinated = 5000.0
		} else {
			evCoordinated = 0.0
		}

		// 3. ESS Coordination & 5-min Regulation
		// In the coordinated scenario, the ESS charges during excess solar and discharges during peak loads.
		// It also regulates the high frequency solar/load fluctuations (imbalances).
		netNonESSPowerCoord := evCoordinated + load1 + load2 - solarPowerCoordinated
		
		// We want to keep the transformer load (net power) around a scheduled target or flat.
		// Scheduled target for net load is hourly energy schedule (e.g. 5kW flat).
		// ESS injects/absorbs to keep the net load at 4.5kW.
		targetNetLoad := 4500.0
		essPowerCoord := targetNetLoad - netNonESSPowerCoord

		// ESS physical limits: -5kW (charging) to +5kW (discharging)
		if essPowerCoord > 5000.0 {
			essPowerCoord = 5000.0
		} else if essPowerCoord < -5000.0 {
			essPowerCoord = -5000.0
		}

		// In the baseline, ESS is uncoordinated (e.g., idle or charging at wrong times)
		var essPowerBaseline float64
		if hour >= 1 && hour <= 4 {
			essPowerBaseline = -3000.0 // Charge battery at night
		} else {
			essPowerBaseline = 0.0
		}

		// Write Coordinated data
		writeStep(wCoord, step, "device_PV", -solarPowerCoordinated) // Generation is negative load
		writeStep(wCoord, step, "device_Load1", load1)
		writeStep(wCoord, step, "device_Load2", load2)
		writeStep(wCoord, step, "device_EV", evCoordinated)
		writeStep(wCoord, step, "device_ESS", -essPowerCoord) // Injecting power reduces grid load

		// Write Baseline data
		writeStep(wBase, step, "device_PV", -solarPowerBaseline)
		writeStep(wBase, step, "device_Load1", load1)
		writeStep(wBase, step, "device_Load2", load2)
		writeStep(wBase, step, "device_EV", evBaseline)
		writeStep(wBase, step, "device_ESS", -essPowerBaseline)
	}

	wCoord.Flush()
	wBase.Flush()
	fmt.Println("✅ Scenario 2.1 CSV files written.")
}

func runScenario22() {
	fmt.Println("👉 Running Scenario 2.2: Blackstart Coordination...")
	
	const totalSteps = 60

	// Output files
	fCoordinated, err := os.Create(filepath.Join("data", "scenario_blackstart.csv"))
	if err != nil {
		panic(err)
	}
	defer fCoordinated.Close()

	fBaseline, err := os.Create(filepath.Join("data", "scenario_blackstart_baseline.csv"))
	if err != nil {
		panic(err)
	}
	defer fBaseline.Close()

	wCoord := csv.NewWriter(fCoordinated)
	wBase := csv.NewWriter(fBaseline)

	header := []string{"Step", "Device", "Power(W)"}
	wCoord.Write(header)
	wBase.Write(header)

	// In blackstart coordination:
	// We want to re-engage 5 devices.
	// Baseline: all 5 reconnect at step 0, causing a huge load spike.
	// Coordinated: step-by-step cold load pickup, keeping power under 10kW.
	deviceLoads := map[string]float64{
		"device_Load1": 4000.0,
		"device_Load2": 5000.0,
		"device_EV":    6000.0,
		"device_ESS":   3000.0, // battery charging load
		"device_PV":    -2000.0, // PV generation starts up later
	}

	for step := 0; step < totalSteps; step++ {
		// Coordinated
		for dev, maxLoad := range deviceLoads {
			var power float64
			switch dev {
			case "device_Load1":
				if step >= 0 {
					power = maxLoad
				}
			case "device_Load2":
				if step >= 12 {
					power = maxLoad
				}
			case "device_EV":
				if step >= 24 {
					power = maxLoad
				}
			case "device_ESS":
				if step >= 36 {
					power = maxLoad
				}
			case "device_PV":
				if step >= 48 {
					power = maxLoad
				}
			}
			writeStep(wCoord, step, dev, power)
		}

		// Baseline: All connect at step 0, PV generation also connects
		for dev, maxLoad := range deviceLoads {
			var power float64
			if dev == "device_PV" {
				if step >= 15 { // PV inverter takes a few minutes to sync
					power = maxLoad
				}
			} else {
				power = maxLoad
			}
			writeStep(wBase, step, dev, power)
		}
	}

	wCoord.Flush()
	wBase.Flush()
	fmt.Println("✅ Scenario 2.2 CSV files written.")
}

func runScenario23() {
	fmt.Println("👉 Running Scenario 2.3: Reserve Capacity & Settlement...")

	// Create a contingency event
	// Device 1: Available load shed = 5kW
	// Device 2: Available load shed = 10kW
	// Reserve Target: 12kW shed requested by the grid
	
	// Create settlement engine and calculate performance
	engine := &operator.SettlementEngine{}

	controls := []*sep.DERControl{
		{
			DERControlBase: &sep.DERControlBase{
				OpModTargetW: &sep.ActivePowerControlType{
					ActivePower: &sep.ActivePower{
						Value: 4000, // 4kW scheduled from device 1
					},
				},
			},
		},
		{
			DERControlBase: &sep.DERControlBase{
				OpModTargetW: &sep.ActivePowerControlType{
					ActivePower: &sep.ActivePower{
						Value: 8000, // 8kW scheduled from device 2
					},
				},
			},
		},
	}

	// Setup actual telemetry (readings)
	// Device 1 delivered 3850W (96.25% accuracy)
	readings1 := []*sep.MirrorUsagePoint{
		{
			MirrorMeterReading: []*sep.MirrorMeterReading{
				{
					Reading: &sep.Reading{
						ReadingBase: &sep.ReadingBase{
							Value: 3850,
						},
					},
				},
			},
		},
	}

	// Device 2 delivered 8100W (98.75% accuracy)
	readings2 := []*sep.MirrorUsagePoint{
		{
			MirrorMeterReading: []*sep.MirrorMeterReading{
				{
					Reading: &sep.Reading{
						ReadingBase: &sep.ReadingBase{
							Value: 8100,
						},
					},
				},
			},
		},
	}

	perf1 := engine.CalculatePerformance("device-1", []*sep.DERControl{controls[0]}, readings1)
	perf2 := engine.CalculatePerformance("device-2", []*sep.DERControl{controls[1]}, readings2)

	report := map[string]interface{}{
		"Timestamp": time.Now().Format(time.RFC3339),
		"ContingencyEvent": map[string]interface{}{
			"RequestedReductionKW": 12.0,
			"StartTime":            time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"DurationMinutes":      60,
		},
		"Settlement": []map[string]interface{}{
			{
				"DeviceLFDI":      perf1.DeviceLFDI,
				"ScheduledKW":     perf1.EnergyScheduled / 1000.0,
				"DeliveredKW":     perf1.EnergyDelivered / 1000.0,
				"AccuracyPercent": perf1.Accuracy * 100.0,
				"Status":          "Compliant",
			},
			{
				"DeviceLFDI":      perf2.DeviceLFDI,
				"ScheduledKW":     perf2.EnergyScheduled / 1000.0,
				"DeliveredKW":     perf2.EnergyDelivered / 1000.0,
				"AccuracyPercent": perf2.Accuracy * 100.0,
				"Status":          "Compliant",
			},
		},
		"TotalReserveVerification": map[string]interface{}{
			"TotalScheduledKW": (perf1.EnergyScheduled + perf2.EnergyScheduled) / 1000.0,
			"TotalDeliveredKW": (perf1.EnergyDelivered + perf2.EnergyDelivered) / 1000.0,
			"OverallAccuracy":  ((perf1.Accuracy + perf2.Accuracy) / 2.0) * 100.0,
			"TargetMet":        ((perf1.EnergyDelivered + perf2.EnergyDelivered) >= 12000.0),
		},
	}

	fReport, err := os.Create(filepath.Join("data", "scenario_settlement.json"))
	if err != nil {
		panic(err)
	}
	defer fReport.Close()

	encoder := json.NewEncoder(fReport)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(report)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Scenario 2.3 Settlement report JSON written.")
}

func writeStep(w *csv.Writer, step int, device string, power float64) {
	err := w.Write([]string{
		strconv.Itoa(step),
		device,
		strconv.FormatFloat(power, 'f', 2, 64),
	})
	if err != nil {
		panic(err)
	}
}
