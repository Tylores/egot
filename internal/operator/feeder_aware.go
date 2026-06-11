package operator

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Tylores/egot/sep"
)

type FeederAwareDispatcher struct {
	scheduler      *GreedyScheduler
	drScheduler    *DRScheduler
	nodeCapacities map[string]float64
	deviceToNode   map[string]string
}

// NewFeederAwareDispatcher creates a new dispatcher with the given scheduler and topology.
func NewFeederAwareDispatcher(scheduler *GreedyScheduler, drScheduler *DRScheduler, topology FeederTopology) *FeederAwareDispatcher {
	caps := make(map[string]float64)
	for _, c := range topology.Capacities {
		caps[c.NodeID] = c.CapacityKW
	}

	maps := make(map[string]string)
	for _, m := range topology.Mappings {
		maps[m.DeviceLFDI] = m.NodeID
	}

	return &FeederAwareDispatcher{
		scheduler:      scheduler,
		drScheduler:    drScheduler,
		nodeCapacities: caps,
		deviceToNode:   maps,
	}
}

// Schedule generates a list of ScheduledEvents while avoiding transformer overloads.
func (d *FeederAwareDispatcher) Schedule(gridReq GridServiceRequest) ([]ScheduledEvent, error) {
	interval := 5 * time.Minute

	// Initialize power balance and node loading for each 5-minute slot
	numSlots := int(gridReq.Duration / interval)
	if gridReq.Duration%interval != 0 {
		numSlots++
	}
	slotBalances := make([]float64, numSlots)
	for i := range slotBalances {
		slotBalances[i] = gridReq.PowerKW
	}

	// nodeSlotLoading[nodeID][slotIndex] = current power in KW
	nodeSlotLoading := make(map[string][]float64)

	// Copy the requests slice to avoid concurrent data races when sorting in-place.
	requests := make([]*sep.FlowReservationRequest, len(d.scheduler.Requests))
	copy(requests, d.scheduler.Requests)

	// Sort requests by absolute power (greedy)
	sort.Slice(requests, func(i, j int) bool {
		pi, pj := 0.0, 0.0
		if requests[i] != nil && requests[i].PowerRequested != nil {
			pi = math.Abs(float64(requests[i].PowerRequested.Value))
		}
		if requests[j] != nil && requests[j].PowerRequested != nil {
			pj = math.Abs(float64(requests[j].PowerRequested.Value))
		}
		return pi > pj
	})

	var scheduled []ScheduledEvent
	usedMRIDs := make(map[string]bool)

	for _, derReq := range requests {
		if derReq.MRID == nil || derReq.MRID.HexBinary128 == nil {
			continue
		}
		mrid := string(*derReq.MRID.HexBinary128)
		if usedMRIDs[mrid] {
			continue
		}
		if derReq.IntervalRequested == nil || derReq.IntervalRequested.Start == nil || derReq.PowerRequested == nil {
			continue
		}

		derStart := time.Unix(int64(*derReq.IntervalRequested.Start), 0)
		derEnd := derStart.Add(time.Duration(derReq.IntervalRequested.Duration) * time.Second)

		// Identify target node
		nodeID := d.deviceToNode[mrid]
		if _, ok := nodeSlotLoading[nodeID]; !ok && nodeID != "" {
			nodeSlotLoading[nodeID] = make([]float64, numSlots)
		}

		// See if this request covers any slots that still need power AND doesn't cause overload
		coversNeededSlot := false
		causesOverload := false
		derPowerKW := float64(derReq.PowerRequested.Value) * 1e-3

		for i := 0; i < numSlots; i++ {
			slotStart := gridReq.StartTime.Add(time.Duration(i) * interval)
			slotEnd := slotStart.Add(interval)

			if (derStart.Before(slotEnd) || derStart.Equal(slotStart)) && derEnd.After(slotStart) {
				// This request overlaps with the slot
				if (gridReq.PowerKW > 0 && slotBalances[i] > 0 && derPowerKW > 0) || (gridReq.PowerKW < 0 && slotBalances[i] < 0 && derPowerKW < 0) {
					coversNeededSlot = true
				}

				// Check feeder capacity at this node for this slot
				if nodeID != "" {
					if limit, ok := d.nodeCapacities[nodeID]; ok {
						if math.Abs(nodeSlotLoading[nodeID][i]+derPowerKW) > limit {
							causesOverload = true
							break
						}
					}
				}
			}
		}

		if coversNeededSlot && !causesOverload {
			// Assign this request
			for i := 0; i < numSlots; i++ {
				slotStart := gridReq.StartTime.Add(time.Duration(i) * interval)
				slotEnd := slotStart.Add(interval)
				if (derStart.Before(slotEnd) || derStart.Equal(slotStart)) && derEnd.After(slotStart) {
					slotBalances[i] -= derPowerKW
					if nodeID != "" {
						nodeSlotLoading[nodeID][i] += derPowerKW
					}
				}
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
								CurrentStatus: 0,
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
		} else if causesOverload {
			fmt.Printf("FeederAwareDispatcher: Overload avoidance triggered for DER %s at node %s\n", mrid, nodeID)
		}
	}

	// Check if all slots are satisfied
	for i, bal := range slotBalances {
		if (gridReq.PowerKW > 0 && bal > 0) || (gridReq.PowerKW < 0 && bal < 0) {
			return scheduled, fmt.Errorf("slot %d unsatisfied: %.2fkW remaining (Feeder constraints may have limited dispatch)", i, bal)
		}
	}

	return scheduled, nil
}

// ScheduleDR generates a list of ScheduledDREvents while avoiding shedding too much load on one feeder segment.
func (d *FeederAwareDispatcher) ScheduleDR(gridReq GridServiceRequest) ([]ScheduledDREvent, error) {
	if d.drScheduler == nil {
		return nil, fmt.Errorf("DR scheduler is not configured")
	}

	// We assume a simple over-voltage risk: don't shed more than 50% of node capacity.
	// In a real scenario, this would be based on minimum load constraints.
	nodeShedLimits := make(map[string]float64)
	for node, cap := range d.nodeCapacities {
		nodeShedLimits[node] = cap * 0.5 // Allow up to 50% of capacity to be shed
	}

	// nodeShedding[nodeID] = amount of load shed scheduled so far in KW
	nodeShedding := make(map[string]float64)

	type deviceAvail struct {
		lfdi string
		lsa  *sep.LoadShedAvailability
	}
	var availList []deviceAvail

	d.drScheduler.mu.RLock()
	for lfdi, lsa := range d.drScheduler.Availability {
		if lsa != nil && lsa.SheddablePower != nil && lsa.SheddablePower.Value > 0 {
			availList = append(availList, deviceAvail{lfdi: lfdi, lsa: lsa})
		}
	}
	d.drScheduler.mu.RUnlock()

	// Sort available devices by SheddablePower (largest first)
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
			break
		}

		if dev.lsa.AvailabilityDuration > 0 && dev.lsa.AvailabilityDuration < durationSec {
			continue
		}

		deviceShedKW := float64(dev.lsa.SheddablePower.Value) * 1e-3
		nodeID := d.deviceToNode[dev.lfdi]

		// Check feeder shed capacity limit
		if nodeID != "" {
			if limit, ok := nodeShedLimits[nodeID]; ok {
				if nodeShedding[nodeID]+deviceShedKW > limit {
					fmt.Printf("FeederAwareDispatcher: Load shed limit reached for node %s\n", nodeID)
					continue
				}
			}
		}

		mrid := "dr-event-" + dev.lfdi[:8]
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
		}

		scheduled = append(scheduled, ScheduledDREvent{
			DeviceLFDI: dev.lfdi,
			Control:    edc,
		})

		shedSoFarKW += deviceShedKW
		if nodeID != "" {
			nodeShedding[nodeID] += deviceShedKW
		}
	}

	if shedSoFarKW < targetShedKW {
		return scheduled, fmt.Errorf("insufficient sheddable power to meet target of %.2fkW (shed %.2fkW)", targetShedKW, shedSoFarKW)
	}

	return scheduled, nil
}
