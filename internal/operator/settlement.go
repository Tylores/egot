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
	totalDelivered := 0.0
	totalScheduled := 0.0

	for _, r := range readings {
		if r == nil {
			continue
		}
		for _, mmr := range r.MirrorMeterReading {
			if mmr != nil && mmr.Reading != nil && mmr.Reading.ReadingBase != nil {
				dt := 3600.0 // default to 1 hour (3600 seconds) if missing
				if mmr.Reading.TimePeriod != nil && mmr.Reading.TimePeriod.Duration > 0 {
					dt = float64(mmr.Reading.TimePeriod.Duration)
				}
				energyKWh := (float64(mmr.Reading.Value) * dt) / (3600.0 * 1000.0)
				totalDelivered += energyKWh
			}
		}
	}

	for _, c := range controls {
		if c != nil && c.DERControlBase != nil && c.DERControlBase.OpModTargetW != nil && c.DERControlBase.OpModTargetW.ActivePower != nil {
			dt := 3600.0 // default to 1 hour if missing
			if c.RandomizableEvent != nil &&
				c.RandomizableEvent.Event != nil &&
				c.RandomizableEvent.Event.Interval != nil &&
				c.RandomizableEvent.Event.Interval.Duration > 0 {
				dt = float64(c.RandomizableEvent.Event.Interval.Duration)
			}
			energyKWh := (float64(c.DERControlBase.OpModTargetW.Value) * dt) / (3600.0 * 1000.0)
			totalScheduled += energyKWh
		}
	}

	accuracy := 1.0
	if totalScheduled != 0 {
		accuracy = 1.0 - math.Abs(totalDelivered-totalScheduled)/math.Abs(totalScheduled)
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
