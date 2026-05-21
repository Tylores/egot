package operator

import (
	"testing"
	"time"

	"github.com/Tylores/egot/sep"
)

func TestGreedyScheduler_Schedule_NonUniform(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	mrid1 := "mrid1"
	mrid2 := "mrid2"
	mrid3 := "mrid3"

	nowUnix := sep.TimeType(now.Unix())
	nowPlus5Unix := sep.TimeType(now.Add(5 * time.Minute).Unix())

	requests := []*sep.FlowReservationRequest{
		{
			IdentifiedObject: &sep.IdentifiedObject{
				MRID: &sep.MRIDType{HexBinary128: &mrid1},
			},
			PowerRequested: &sep.ActivePower{Value: 10000}, // 10kW
			IntervalRequested: &sep.DateTimeInterval{
				Start:    &nowUnix,
				Duration: 300, // 5 min
			},
		},
		{
			IdentifiedObject: &sep.IdentifiedObject{
				MRID: &sep.MRIDType{HexBinary128: &mrid2},
			},
			PowerRequested: &sep.ActivePower{Value: 5000}, // 5kW
			IntervalRequested: &sep.DateTimeInterval{
				Start:    &nowUnix,
				Duration: 900, // 15 min
			},
		},
		{
			IdentifiedObject: &sep.IdentifiedObject{
				MRID: &sep.MRIDType{HexBinary128: &mrid3},
			},
			PowerRequested: &sep.ActivePower{Value: 10000}, // 10kW
			IntervalRequested: &sep.DateTimeInterval{
				Start:    &nowPlus5Unix,
				Duration: 300, // 5 min
			},
		},
	}

	scheduler := NewScheduler(requests)

	// Grid request: 15kW for 10 minutes starting now
	gridReq := GridServiceRequest{
		StartTime: now,
		Duration:  10 * time.Minute,
		PowerKW:   15.0,
	}

	scheduled, err := scheduler.Schedule(gridReq)
	if err != nil {
		t.Fatalf("Schedule failed: %v", err)
	}

	// Expected:
	// Slot 1 (0-5 min): 15kW needed. mrid1 (10kW) + mrid2 (5kW) = 15kW.
	// Slot 2 (5-10 min): 15kW needed. mrid2 (5kW) + mrid3 (10kW) = 15kW.
	// All three should be scheduled.
	if len(scheduled) != 3 {
		t.Errorf("Expected 3 scheduled events, got %d", len(scheduled))
	}

	// Check if all MRIDs are present
	mridSet := make(map[string]bool)
	for _, ev := range scheduled {
		mridSet[string(*ev.Control.MRID.HexBinary128)] = true
	}

	if !mridSet[mrid1] || !mridSet[mrid2] || !mridSet[mrid3] {
		t.Errorf("One or more MRIDs missing from scheduled events: %v", mridSet)
	}
}

func TestGreedyScheduler_Schedule_Unsatisfied(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	nowUnix := sep.TimeType(now.Unix())
	mrid1 := "mrid1"

	requests := []*sep.FlowReservationRequest{
		{
			IdentifiedObject: &sep.IdentifiedObject{
				MRID: &sep.MRIDType{HexBinary128: &mrid1},
			},
			PowerRequested: &sep.ActivePower{Value: 10000}, // 10kW
			IntervalRequested: &sep.DateTimeInterval{
				Start:    &nowUnix,
				Duration: 300, // 5 min
			},
		},
	}


	scheduler := NewScheduler(requests)

	// Grid request: 15kW for 5 minutes
	gridReq := GridServiceRequest{
		StartTime: now,
		Duration:  5 * time.Minute,
		PowerKW:   15.0,
	}

	_, err := scheduler.Schedule(gridReq)
	if err == nil {
		t.Errorf("Expected error for unsatisfied request, got nil")
	}
}

func TestGreedyScheduler_DayAhead(t *testing.T) {
	startOfDay := time.Now().Truncate(24 * time.Hour)

	// Three DERs with requests spanning different parts of the 24-hour window.
	mrid1 := "da-mrid-1" // Available 8am–6pm (10 hours)
	mrid2 := "da-mrid-2" // Available 6pm–10pm (4 hours)
	mrid3 := "da-mrid-3" // Available midnight–8am (8 hours)

	t8am := sep.TimeType(startOfDay.Add(8 * time.Hour).Unix())
	t6pm := sep.TimeType(startOfDay.Add(18 * time.Hour).Unix())
	tmid := sep.TimeType(startOfDay.Unix())

	requests := []*sep.FlowReservationRequest{
		{
			IdentifiedObject:  &sep.IdentifiedObject{MRID: &sep.MRIDType{HexBinary128: &mrid1}},
			PowerRequested:    &sep.ActivePower{Value: 15000}, // 15kW
			IntervalRequested: &sep.DateTimeInterval{Start: &t8am, Duration: 36000}, // 10h
		},
		{
			IdentifiedObject:  &sep.IdentifiedObject{MRID: &sep.MRIDType{HexBinary128: &mrid2}},
			PowerRequested:    &sep.ActivePower{Value: 8000}, // 8kW
			IntervalRequested: &sep.DateTimeInterval{Start: &t6pm, Duration: 14400}, // 4h
		},
		{
			IdentifiedObject:  &sep.IdentifiedObject{MRID: &sep.MRIDType{HexBinary128: &mrid3}},
			PowerRequested:    &sep.ActivePower{Value: 5000}, // 5kW
			IntervalRequested: &sep.DateTimeInterval{Start: &tmid, Duration: 28800}, // 8h
		},
	}

	scheduler := NewScheduler(requests)
	scheduled, err := scheduler.ScheduleDayAhead(startOfDay)
	if err != nil {
		t.Fatalf("ScheduleDayAhead failed: %v", err)
	}

	if len(scheduled) != 3 {
		t.Errorf("Expected 3 scheduled events, got %d", len(scheduled))
	}

	// Verify all three MRIDs are represented.
	mridSet := make(map[string]bool)
	for _, ev := range scheduled {
		mridSet[string(*ev.Control.MRID.HexBinary128)] = true
	}
	for _, m := range []string{mrid1, mrid2, mrid3} {
		if !mridSet[m] {
			t.Errorf("Expected MRID %s in day-ahead schedule, not found", m)
		}
	}

	t.Logf("Day-ahead schedule: %d events across 24h window", len(scheduled))
}

func TestGreedyScheduler_EIMWindow(t *testing.T) {
	now := time.Now().Truncate(time.Hour)
	mrid1 := "eim-mrid-1"
	mrid2 := "eim-mrid-2"

	nowUnix := sep.TimeType(now.Unix())

	requests := []*sep.FlowReservationRequest{
		{
			IdentifiedObject:  &sep.IdentifiedObject{MRID: &sep.MRIDType{HexBinary128: &mrid1}},
			PowerRequested:    &sep.ActivePower{Value: 8000}, // 8kW
			IntervalRequested: &sep.DateTimeInterval{Start: &nowUnix, Duration: 300}, // 5 min
		},
		{
			IdentifiedObject:  &sep.IdentifiedObject{MRID: &sep.MRIDType{HexBinary128: &mrid2}},
			PowerRequested:    &sep.ActivePower{Value: 7000}, // 7kW
			IntervalRequested: &sep.DateTimeInterval{Start: &nowUnix, Duration: 300}, // 5 min
		},
	}

	scheduler := NewScheduler(requests)

	// EIM: need 12kW for 5 minutes
	gridReq := GridServiceRequest{
		StartTime: now,
		Duration:  5 * time.Minute,
		PowerKW:   12.0,
	}

	scheduled, err := scheduler.ScheduleWithWindow(gridReq, EIMWindow)
	if err != nil {
		t.Fatalf("EIM schedule failed: %v", err)
	}

	totalPower := 0.0
	for _, ev := range scheduled {
		totalPower += float64(ev.Control.DERControlBase.OpModTargetW.Value) * 1e-3
	}

	if totalPower < 12.0 {
		t.Errorf("Expected at least 12.0kW scheduled, got %.2fkW", totalPower)
	}

	t.Logf("EIM schedule: %.2fkW across %d devices", totalPower, len(scheduled))
}

