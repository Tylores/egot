package operator

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Tylores/egot/sep"
)

// ScheduleWindow defines the scheduling resolution window.
type ScheduleWindow int

const (
	// EIMWindow schedules in 5-minute slots for Energy Imbalance Market real-time balancing.
	EIMWindow ScheduleWindow = iota
	// DayAheadWindow schedules in 60-minute (hourly) slots for 24-hour day-ahead planning.
	DayAheadWindow
)

// slotDuration returns the interval size for the given window type.
func (w ScheduleWindow) slotDuration() time.Duration {
	switch w {
	case DayAheadWindow:
		return time.Hour
	default:
		return 5 * time.Minute
	}
}

// GridServiceRequest represents a demand from the grid (e.g., "we need 100kW reduction at 4PM").
type GridServiceRequest struct {
	StartTime time.Time
	Duration  time.Duration
	PowerKW   float64 // Positive for injection, negative for load reduction
}

// ScheduledEvent represents a DERControl that has been assigned to a DER.
type ScheduledEvent struct {
	DeviceLFDI string
	Control    *sep.DERControl
}

// GreedyScheduler picks DERs to satisfy a GridServiceRequest.
type GreedyScheduler struct {
	Requests []*sep.FlowReservationRequest
}

// NewScheduler creates a new GreedyScheduler with the given available requests.
func NewScheduler(requests []*sep.FlowReservationRequest) *GreedyScheduler {
	return &GreedyScheduler{Requests: requests}
}

// Schedule generates a list of ScheduledEvents using the EIM (5-minute) window.
// This is a convenience wrapper for ScheduleWithWindow(gridReq, EIMWindow).
func (s *GreedyScheduler) Schedule(gridReq GridServiceRequest) ([]ScheduledEvent, error) {
	return s.ScheduleWithWindow(gridReq, EIMWindow)
}

// ScheduleDayAhead generates a 24-hour, hourly-slotted schedule for day-ahead planning.
// It expects requests that span multiple hours.
func (s *GreedyScheduler) ScheduleDayAhead(startOfDay time.Time) ([]ScheduledEvent, error) {
	gridReq := GridServiceRequest{
		StartTime: startOfDay.Truncate(time.Hour),
		Duration:  24 * time.Hour,
		PowerKW:   0, // Not a single power target — individual hourly slot balancing
	}
	return s.scheduleHourly(gridReq)
}

// scheduleHourly runs the day-ahead algorithm: for each hourly slot, greedily assign
// available DER requests to maximize utilization across the 24-hour window.
func (s *GreedyScheduler) scheduleHourly(gridReq GridServiceRequest) ([]ScheduledEvent, error) {
	interval := DayAheadWindow.slotDuration()
	numSlots := int(gridReq.Duration / interval)

	// slotPowerKW[i] tracks the net power committed for slot i (positive = injection).
	slotPowerKW := make([]float64, numSlots)

	// Sort requests by absolute power descending to maximise per-assignment impact.
	sort.Slice(s.Requests, func(i, j int) bool {
		pi, pj := 0.0, 0.0
		if s.Requests[i].PowerRequested != nil {
			pi = math.Abs(float64(s.Requests[i].PowerRequested.Value))
		}
		if s.Requests[j].PowerRequested != nil {
			pj = math.Abs(float64(s.Requests[j].PowerRequested.Value))
		}
		return pi > pj
	})

	var scheduled []ScheduledEvent
	usedMRIDs := make(map[string]bool)

	for _, derReq := range s.Requests {
		if derReq.MRID == nil || derReq.MRID.HexBinary128 == nil {
			continue
		}
		mrid := string(*derReq.MRID.HexBinary128)
		if usedMRIDs[mrid] {
			continue
		}

		if derReq.IntervalRequested == nil || derReq.IntervalRequested.Start == nil {
			continue
		}
		if derReq.PowerRequested == nil {
			continue
		}

		derStart := time.Unix(int64(*derReq.IntervalRequested.Start), 0)
		derEnd := derStart.Add(time.Duration(derReq.IntervalRequested.Duration) * time.Second)
		derPowerKW := float64(derReq.PowerRequested.Value) * 1e-3

		// Find which hourly slots this request overlaps with.
		var overlappingSlots []int
		for i := 0; i < numSlots; i++ {
			slotStart := gridReq.StartTime.Add(time.Duration(i) * interval)
			slotEnd := slotStart.Add(interval)
			if (derStart.Before(slotEnd) || derStart.Equal(slotStart)) && derEnd.After(slotStart) {
				overlappingSlots = append(overlappingSlots, i)
			}
		}

		if len(overlappingSlots) == 0 {
			continue // Request falls outside the day-ahead window
		}

		// Commit this DER to all overlapping slots.
		for _, i := range overlappingSlots {
			slotPowerKW[i] += derPowerKW
		}

		scheduled = append(scheduled, ScheduledEvent{
			DeviceLFDI: mrid,
			Control: &sep.DERControl{
				RandomizableEvent: &sep.RandomizableEvent{
					Event: &sep.Event{
						RespondableSubscribableIdentifiedObject: &sep.RespondableSubscribableIdentifiedObject{
							MRID: &sep.MRIDType{HexBinary128: derReq.MRID.HexBinary128},
						},
						EventStatus: &sep.EventStatus{
							CurrentStatus: 0, // Scheduled
						},
						Interval: derReq.IntervalRequested,
					},
				},
				DERControlBase: &sep.DERControlBase{
					OpModTargetW: &sep.ActivePowerControlType{
						ActivePower: &sep.ActivePower{
							Value: int16(derReq.PowerRequested.Value),
						},
					},
				},
			},
		})
		usedMRIDs[mrid] = true
	}

	return scheduled, nil
}

// ScheduleWithWindow generates ScheduledEvents using the specified window resolution.
func (s *GreedyScheduler) ScheduleWithWindow(gridReq GridServiceRequest, window ScheduleWindow) ([]ScheduledEvent, error) {
	interval := window.slotDuration()

	// Initialize power balance for each slot
	numSlots := int(gridReq.Duration / interval)
	if gridReq.Duration%interval != 0 {
		numSlots++
	}
	slotBalances := make([]float64, numSlots)
	for i := range slotBalances {
		slotBalances[i] = gridReq.PowerKW
	}

	// Sort requests by absolute power (greedy)
	sort.Slice(s.Requests, func(i, j int) bool {
		pi, pj := 0.0, 0.0
		if s.Requests[i].PowerRequested != nil {
			pi = math.Abs(float64(s.Requests[i].PowerRequested.Value))
		}
		if s.Requests[j].PowerRequested != nil {
			pj = math.Abs(float64(s.Requests[j].PowerRequested.Value))
		}
		return pi > pj
	})

	var scheduled []ScheduledEvent
	usedMRIDs := make(map[string]bool)

	for _, derReq := range s.Requests {
		if derReq.MRID == nil || derReq.MRID.HexBinary128 == nil {
			continue
		}
		mrid := string(*derReq.MRID.HexBinary128)
		if usedMRIDs[mrid] {
			continue
		}

		if derReq.IntervalRequested == nil || derReq.IntervalRequested.Start == nil {
			continue
		}
		derStart := time.Unix(int64(*derReq.IntervalRequested.Start), 0)
		derEnd := derStart.Add(time.Duration(derReq.IntervalRequested.Duration) * time.Second)

		// See if this request covers any slots that still need power
		coversNeededSlot := false
		derPowerKW := float64(derReq.PowerRequested.Value) * 1e-3
		for i := 0; i < numSlots; i++ {
			slotStart := gridReq.StartTime.Add(time.Duration(i) * interval)
			slotEnd := slotStart.Add(interval)

			if (derStart.Before(slotEnd) || derStart.Equal(slotStart)) && derEnd.After(slotStart) {
				// This request overlaps with the slot
				if (gridReq.PowerKW > 0 && slotBalances[i] > 0 && derPowerKW > 0) || (gridReq.PowerKW < 0 && slotBalances[i] < 0 && derPowerKW < 0) {
					coversNeededSlot = true
					break
				}
			}
		}

		if coversNeededSlot {
			// Assign this request
			for i := 0; i < numSlots; i++ {
				slotStart := gridReq.StartTime.Add(time.Duration(i) * interval)
				slotEnd := slotStart.Add(interval)
				if (derStart.Before(slotEnd) || derStart.Equal(slotStart)) && derEnd.After(slotStart) {
					slotBalances[i] -= derPowerKW
				}
			}

			scheduled = append(scheduled, ScheduledEvent{
				DeviceLFDI: mrid, // Using mRID as proxy for LFDI
				Control: &sep.DERControl{
					RandomizableEvent: &sep.RandomizableEvent{
						Event: &sep.Event{
							RespondableSubscribableIdentifiedObject: &sep.RespondableSubscribableIdentifiedObject{
								MRID: &sep.MRIDType{HexBinary128: derReq.MRID.HexBinary128},
							},
							EventStatus: &sep.EventStatus{
								CurrentStatus: 0, // Scheduled
							},
							Interval: derReq.IntervalRequested,
						},
					},
					DERControlBase: &sep.DERControlBase{
						OpModTargetW: &sep.ActivePowerControlType{
							ActivePower: &sep.ActivePower{
								Value: int16(derReq.PowerRequested.Value),
							},
						},
					},
				},
			})
			usedMRIDs[mrid] = true
		}
	}

	// Check if all slots are satisfied
	for i, bal := range slotBalances {
		if (gridReq.PowerKW > 0 && bal > 0) || (gridReq.PowerKW < 0 && bal < 0) {
			return scheduled, fmt.Errorf("slot %d (start %v) unsatisfied: %.2fkW remaining", i, gridReq.StartTime.Add(time.Duration(i)*interval), bal)
		}
	}

	return scheduled, nil
}
