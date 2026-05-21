package operator

import (
	"testing"
	"time"

	"github.com/Tylores/egot/sep"
)

func TestFeederAwareDispatcher_OverloadAvoidance(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	nowUnix := sep.TimeType(now.Unix())
	mrid1 := "mrid1"
	mrid2 := "mrid2"
	nodeA := "nodeA"

	// Two DERs on the same node, each 10kW
	requests := []*sep.FlowReservationRequest{
		{
			IdentifiedObject: &sep.IdentifiedObject{
				MRID: &sep.MRIDType{HexBinary128: &mrid1},
			},
			PowerRequested: &sep.ActivePower{Value: 10000},
			IntervalRequested: &sep.DateTimeInterval{
				Start:    &nowUnix,
				Duration: 300,
			},
		},
		{
			IdentifiedObject: &sep.IdentifiedObject{
				MRID: &sep.MRIDType{HexBinary128: &mrid2},
			},
			PowerRequested: &sep.ActivePower{Value: 10000},
			IntervalRequested: &sep.DateTimeInterval{
				Start:    &nowUnix,
				Duration: 300,
			},
		},
	}

	topology := FeederTopology{
		Capacities: []NodeCapacity{
			{NodeID: nodeA, CapacityKW: 15.0}, // Only 15kW allowed at nodeA
		},
		Mappings: []DeviceMapping{
			{DeviceLFDI: mrid1, NodeID: nodeA},
			{DeviceLFDI: mrid2, NodeID: nodeA},
		},
	}

	scheduler := NewScheduler(requests)
	dispatcher := NewFeederAwareDispatcher(scheduler, nil, topology)

	// Grid request: 20kW needed
	gridReq := GridServiceRequest{
		StartTime: now,
		Duration:  5 * time.Minute,
		PowerKW:   20.0,
	}

	scheduled, err := dispatcher.Schedule(gridReq)
	
	// We expect an error because 20kW cannot be satisfied without overloading nodeA
	if err == nil {
		t.Errorf("Expected error due to unsatisfied request (feeder limit), got nil")
	}

	// But we should still have 1 DER scheduled (the 10kW one)
	if len(scheduled) != 1 {
		t.Errorf("Expected 1 DER to be scheduled within feeder limits, got %d", len(scheduled))
	}
}
