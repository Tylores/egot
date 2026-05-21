package operator

import (
	"math"

	"github.com/Tylores/egot/sep"
)

// PerformanceReport summarizes how well a DER followed its instructions.
type PerformanceReport struct {
	DeviceLFDI       string
	Accuracy         float64 // 0.0 to 1.0
	EnergyDelivered  float64 // kWh
	EnergyScheduled  float64 // kWh
}

// SettlementEngine compares schedules against actual telemetry.
type SettlementEngine struct{}

// CalculatePerformance compares a list of controls against usage telemetry.
func (e *SettlementEngine) CalculatePerformance(lfdi string, controls []*sep.DERControl, readings []*sep.MirrorUsagePoint) PerformanceReport {
	// Simple integration of readings to find energy delivered
	// And comparison against the scheduled controls.
	
	totalDelivered := 0.0
	totalScheduled := 0.0
	
	// This is a placeholder for actual performance math.
	// In a real PhD implementation, you'd use Root Mean Square Error (RMSE)
	// or similar metrics to track setpoint following.
	
	for _, r := range readings {
		for _, mmr := range r.MirrorMeterReading {
			if mmr.Reading != nil {
				totalDelivered += float64(mmr.Reading.Value)
			}
		}
	}
	
	for _, c := range controls {
		if c.DERControlBase != nil && c.DERControlBase.OpModTargetW != nil {
			totalScheduled += float64(c.DERControlBase.OpModTargetW.Value)
		}
	}
	
	accuracy := 1.0
	if totalScheduled > 0 {
		accuracy = 1.0 - math.Abs(totalDelivered-totalScheduled)/totalScheduled
	}
	if accuracy < 0 {
		accuracy = 0
	}

	return PerformanceReport{
		DeviceLFDI:      lfdi,
		Accuracy:        accuracy,
		EnergyDelivered: totalDelivered,
		EnergyScheduled: totalScheduled,
	}
}
