package integration

import (
	"testing"
	"time"

	"github.com/Tylores/egot/internal/operator"
	"github.com/Tylores/egot/sep"
)

func TestBlackstartDR_Basic(t *testing.T) {
	now := time.Now().Truncate(time.Minute)

	// Emulate LoadShedAvailability for a few devices
	availMap := map[string]*sep.LoadShedAvailability{
		"device-1": {
			SheddablePower:       &sep.ActivePower{Value: 5000}, // 5kW
			AvailabilityDuration: 3600,                          // 1 hour
		},
		"device-2": {
			SheddablePower:       &sep.ActivePower{Value: 10000}, // 10kW
			AvailabilityDuration: 3600,
		},
		"device-3": {
			SheddablePower:       &sep.ActivePower{Value: 2000}, // 2kW
			AvailabilityDuration: 1800,                          // 30 min
		},
	}

	drScheduler := operator.NewDRScheduler(availMap)

	// Blackstart Grid Request: We need to shed 12kW of load immediately for 45 minutes
	gridReq := operator.GridServiceRequest{
		StartTime: now,
		Duration:  45 * time.Minute,
		PowerKW:   12.0, // We need to shed 12kW
	}

	scheduled, err := drScheduler.ScheduleDR(gridReq)
	if err != nil {
		t.Fatalf("Failed to schedule DR: %v", err)
	}

	totalShed := 0.0
	for _, ev := range scheduled {
		totalShed += float64(availMap[ev.DeviceLFDI].SheddablePower.Value) * 1e-3
	}

	// We expect device-2 (10kW) and device-1 (5kW) to be scheduled, providing 15kW shed.
	// device-3 (2kW) cannot be scheduled because its AvailabilityDuration (30m) is less than the request (45m).
	if totalShed < 12.0 {
		t.Errorf("Expected at least 12kW shed, got %.2fkW", totalShed)
	}

	t.Logf("Successfully scheduled %.2fkW of load shed", totalShed)
}

func TestBlackstartDR_FeederAware(t *testing.T) {
	now := time.Now().Truncate(time.Minute)

	availMap := map[string]*sep.LoadShedAvailability{
		"device-1": {SheddablePower: &sep.ActivePower{Value: 10000}, AvailabilityDuration: 3600},
		"device-2": {SheddablePower: &sep.ActivePower{Value: 10000}, AvailabilityDuration: 3600},
	}

	drScheduler := operator.NewDRScheduler(availMap)
	dummyScheduler := operator.NewScheduler(nil)

	// Topology limits:
	// Node A capacity is 15kW. The dispatcher limits DR shed to 50% of node capacity = 7.5kW.
	topology := operator.FeederTopology{
		Capacities: []operator.NodeCapacity{
			{NodeID: "NodeA", CapacityKW: 15.0},
		},
		Mappings: []operator.DeviceMapping{
			{DeviceLFDI: "device-1", NodeID: "NodeA"},
			{DeviceLFDI: "device-2", NodeID: "NodeA"},
		},
	}

	dispatcher := operator.NewFeederAwareDispatcher(dummyScheduler, drScheduler, topology)

	gridReq := operator.GridServiceRequest{
		StartTime: now,
		Duration:  1 * time.Hour,
		PowerKW:   15.0, // Need to shed 15kW
	}

	scheduled, err := dispatcher.ScheduleDR(gridReq)

	// We expect this to fail or limit the shed because 2 devices in Node A each have 10kW shed.
	// The node limit is 7.5kW (50% of 15kW). Neither device can be fully shed without exceeding the limit.
	// Wait, the logic in FeederAwareDispatcher says: if nodeShedding[nodeID] + deviceShedKW > limit, continue.
	// Since limit is 7.5kW and device shed is 10kW, it will never shed any device on this node!

	if err == nil {
		t.Errorf("Expected failure due to node shed limits, but got success. Scheduled: %v", scheduled)
	} else {
		t.Logf("Correctly rejected due to limits: %v", err)
	}
}
