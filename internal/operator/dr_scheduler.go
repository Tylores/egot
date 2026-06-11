package operator

import (
	"math"
	"sort"
	"sync"

	"github.com/Tylores/egot/sep"
)

// ScheduledDREvent represents a Demand Response event (EndDeviceControl) that has been assigned to a device.
type ScheduledDREvent struct {
	DeviceLFDI string
	Control    *sep.EndDeviceControl
}

// DRScheduler picks devices to shed load to satisfy a GridServiceRequest during Blackstart or emergencies.
type DRScheduler struct {
	mu           sync.RWMutex
	Availability map[string]*sep.LoadShedAvailability // Map of LFDI to their availability
}

// NewDRScheduler creates a new DRScheduler with the given available load shed resources.
func NewDRScheduler(availability map[string]*sep.LoadShedAvailability) *DRScheduler {
	return &DRScheduler{Availability: availability}
}

// UpdateAvailability updates the availability map in a thread-safe manner.
func (s *DRScheduler) UpdateAvailability(lfdi string, lsa *sep.LoadShedAvailability) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Availability == nil {
		s.Availability = make(map[string]*sep.LoadShedAvailability)
	}
	s.Availability[lfdi] = lsa
}

// ScheduleDR generates a list of ScheduledDREvents to shed load to meet the grid request.
// The gridReq.PowerKW here represents the total amount of load (in kW) that needs to be *shed* (removed from the grid).
func (s *DRScheduler) ScheduleDR(gridReq GridServiceRequest) ([]ScheduledDREvent, error) {
	// Sort available devices by SheddablePower (largest first) to minimize the number of control signals sent
	type deviceAvail struct {
		lfdi string
		lsa  *sep.LoadShedAvailability
	}
	var availList []deviceAvail

	s.mu.RLock()
	for lfdi, lsa := range s.Availability {
		if lsa != nil && lsa.SheddablePower != nil && lsa.SheddablePower.Value > 0 {
			availList = append(availList, deviceAvail{lfdi: lfdi, lsa: lsa})
		}
	}
	s.mu.RUnlock()

	sort.Slice(availList, func(i, j int) bool {
		return availList[i].lsa.SheddablePower.Value > availList[j].lsa.SheddablePower.Value
	})

	var scheduled []ScheduledDREvent
	targetShedKW := math.Abs(gridReq.PowerKW)
	shedSoFarKW := 0.0

	nowUnix := sep.TimeType(gridReq.StartTime.Unix())
	durationSec := uint32(gridReq.Duration.Seconds())

	for _, dev := range availList {
		if shedSoFarKW >= targetShedKW {
			break // We have shed enough load
		}

		// Check if the device is available for the required duration
		if dev.lsa.AvailabilityDuration > 0 && dev.lsa.AvailabilityDuration < durationSec {
			continue // Device cannot shed load for the full requested duration
		}

		deviceShedKW := float64(dev.lsa.SheddablePower.Value) * 1e-3

		// We do a simple greedy full-shed for blackstart (no partial shedding via duty cycle for now)
		// Send EndDeviceControl to tell it to reduce load.
		mrid := "dr-event-" + dev.lfdi[:8] // Shortened MRID for simplicity

		edc := &sep.EndDeviceControl{
			RandomizableEvent: &sep.RandomizableEvent{
				Event: &sep.Event{
					RespondableSubscribableIdentifiedObject: &sep.RespondableSubscribableIdentifiedObject{
						MRID: &sep.MRIDType{HexBinary128: &mrid},
					},
					Interval: &sep.DateTimeInterval{
						Start:    &nowUnix,
						Duration: durationSec,
					},
				},
			},
			// ApplianceLoadReduction can be added if needed, or simply let the event interval dictate the shed time
		}

		scheduled = append(scheduled, ScheduledDREvent{
			DeviceLFDI: dev.lfdi,
			Control:    edc,
		})

		shedSoFarKW += deviceShedKW
	}

	return scheduled, nil
}
