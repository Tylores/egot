package integration

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	derHandler "github.com/Tylores/egot/internal/DER/handler"
	drHandler "github.com/Tylores/egot/internal/DR/handler"
	edevHandler "github.com/Tylores/egot/internal/EDevice/handler"
	frHandler "github.com/Tylores/egot/internal/FlowReservation/handler"
	"github.com/Tylores/egot/internal/operator"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/sep"
)

func TestGMLCGridServices_FullFlow(t *testing.T) {
	// 1. Setup Shared Infrastructure
	repo := store.New("test_integration.db")
	repo.Load()
	defer repo.Close()
	reg := registry.New("test_integration_reg.db")
	reg.Load()

	// 2. Initialize Handlers
	hEDev := edevHandler.NewHandler(repo, reg)
	_ = hEDev
	hDER := derHandler.NewHandler(repo, reg)
	hDR := drHandler.NewHandler(repo, reg)
	hFR := frHandler.NewHandler(repo, reg)

	// 3. Mock Device Onboarding
	cert := &x509.Certificate{Raw: []byte("test-integration-cert")}
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	sfdi, _ := sep.ToSFDI(lfdi)
	reg.RegisterLFDI(lfdi, fmt.Sprintf("%d", sfdi), "Integration Test Device")

	tlsState := &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}

	// 4. Grid Service Emulation: Load Shed Availability
	lsa := &sep.LoadShedAvailability{
		SheddablePower:       &sep.ActivePower{Value: 15000},
		AvailabilityDuration: 3600,
	}
	bodyLSA, _ := xml.Marshal(lsa)
	reqLSA := httptest.NewRequest("POST", "/edev/1/lsl", bytes.NewBuffer(bodyLSA))
	reqLSA.TLS = tlsState
	reqLSA.SetPathValue("id1", "1")
	wLSA := httptest.NewRecorder()
	hEDev.POSTLoadShedAvailabilityList(wLSA, reqLSA)
	if wLSA.Code != http.StatusCreated {
		t.Fatalf("Failed to post LoadShedAvailability: %d", wLSA.Code)
	}

	// Verify LSA retrieval
	reqGSAL := httptest.NewRequest("GET", "/edev/1/lsl", nil)
	reqGSAL.TLS = tlsState
	reqGSAL.SetPathValue("id1", "1")
	wGSAL := httptest.NewRecorder()
	hEDev.GETLoadShedAvailabilityList(wGSAL, reqGSAL)
	var lsaList sep.LoadShedAvailabilityList
	xml.NewDecoder(wGSAL.Body).Decode(&lsaList)
	if len(lsaList.LoadShedAvailability) == 0 || lsaList.LoadShedAvailability[0].SheddablePower.Value != 15000 {
		t.Fatalf("LoadShedAvailability not retrieved correctly: %+v", lsaList)
	}
	
	// 5. Grid Service Emulation: DER Curves (Volt-Var)
	// ... (rest of the code remains similar)
	curveMRID := "volt-var-curve-1"
	vvCurve := &sep.DERCurve{
		IdentifiedObject: &sep.IdentifiedObject{
			MRID: &sep.MRIDType{HexBinary128: &curveMRID},
		},
		CurveType: func() *sep.DERCurveType { v := sep.DERCurveType(11); return &v }(),
		CurveData: []*sep.CurveData{
			{Xvalue: 90, Yvalue: 100},
			{Xvalue: 110, Yvalue: -100},
		},
	}
	bodyCurve, _ := xml.Marshal(vvCurve)
	reqCurve := httptest.NewRequest("POST", "/derp/1/dc", bytes.NewBuffer(bodyCurve))
	reqCurve.TLS = tlsState
	reqCurve.SetPathValue("id1", "1")
	wCurve := httptest.NewRecorder()
	hDER.POSTDERCurveList(wCurve, reqCurve)
	if wCurve.Code != http.StatusCreated {
		t.Fatalf("Failed to post DER curve: %d", wCurve.Code)
	}

	// 6. Demand Response Signaling
	now := time.Now().Truncate(time.Minute)
	nowUnix := sep.TimeType(now.Unix())

	edcMRID := "dr-signal-001"
	edc := &sep.EndDeviceControl{
		RandomizableEvent: &sep.RandomizableEvent{
			Event: &sep.Event{
				RespondableSubscribableIdentifiedObject: &sep.RespondableSubscribableIdentifiedObject{
					MRID: &sep.MRIDType{HexBinary128: &edcMRID},
				},
				Interval: &sep.DateTimeInterval{
					Start:    &nowUnix,
					Duration: 3600, // 1 hour DR event
				},
			},
		},
	}
	bodyEDC, _ := xml.Marshal(edc)
	reqEDC := httptest.NewRequest("POST", "/dr/1/edc", bytes.NewBuffer(bodyEDC))
	reqEDC.TLS = tlsState
	reqEDC.SetPathValue("id1", "1")
	wEDC := httptest.NewRecorder()
	hDR.POSTEndDeviceControlList(wEDC, reqEDC)
	if wEDC.Code != http.StatusCreated {
		t.Fatalf("Failed to post DR signal: %d", wEDC.Code)
	}

	// 7. Client Response: Flexible Flow Reservations

	// FRQ 1: 5 minutes, 10kW
	mrid1 := "frq-001"
	frq1 := &sep.FlowReservationRequest{
		IdentifiedObject: &sep.IdentifiedObject{
			MRID: &sep.MRIDType{HexBinary128: &mrid1},
		},
		PowerRequested: &sep.ActivePower{Value: 10000},
		IntervalRequested: &sep.DateTimeInterval{
			Start:    &nowUnix,
			Duration: 300,
		},
		DurationRequested: 300, // Valid EIM duration
	}
	bodyFRQ1, _ := xml.Marshal(frq1)
	reqFRQ1 := httptest.NewRequest("POST", "/frq", bytes.NewBuffer(bodyFRQ1))
	reqFRQ1.TLS = tlsState
	wFRQ1 := httptest.NewRecorder()
	hFR.POSTFlowReservationRequestList(wFRQ1, reqFRQ1)
	if wFRQ1.Code != http.StatusCreated {
		t.Fatalf("Failed to post FRQ1: %d", wFRQ1.Code)
	}

	// Test Invalid Duration (Task 1.1)
	frqInvalid := &sep.FlowReservationRequest{
		DurationRequested: 100, // Invalid
	}
	bodyInvalid, _ := xml.Marshal(frqInvalid)
	reqInvalid := httptest.NewRequest("POST", "/frq", bytes.NewBuffer(bodyInvalid))
	reqInvalid.TLS = tlsState
	wInvalid := httptest.NewRecorder()
	hFR.POSTFlowReservationRequestList(wInvalid, reqInvalid)
	if wInvalid.Code != http.StatusBadRequest {
		t.Fatalf("Expected BadRequest for invalid duration, got %d", wInvalid.Code)
	}

	// FRQ 2: 15 minutes, 5kW
	mrid2 := "frq-002"
	frq2 := &sep.FlowReservationRequest{
		IdentifiedObject: &sep.IdentifiedObject{
			MRID: &sep.MRIDType{HexBinary128: &mrid2},
		},
		PowerRequested: &sep.ActivePower{Value: 5000},
		IntervalRequested: &sep.DateTimeInterval{
			Start:    &nowUnix,
			Duration: 900,
		},
		DurationRequested: 900, // Valid EIM duration
	}
	bodyFRQ2, _ := xml.Marshal(frq2)
	reqFRQ2 := httptest.NewRequest("POST", "/frq", bytes.NewBuffer(bodyFRQ2))
	reqFRQ2.TLS = tlsState
	wFRQ2 := httptest.NewRecorder()
	hFR.POSTFlowReservationRequestList(wFRQ2, reqFRQ2)
	if wFRQ2.Code != http.StatusCreated {
		t.Fatalf("Failed to post FRQ2: %d", wFRQ2.Code)
	}

	// 8. Operator Dispatch with Feeder-Awareness
	// Fetch all requests
	var allRequests []*sep.FlowReservationRequest
	results := repo.GetByOwner(fmt.Sprintf("%d", sfdi))
	for _, res := range results {
		if frq, ok := res.(*sep.FlowReservationRequest); ok {
			allRequests = append(allRequests, frq)
		}
	}

	if len(allRequests) < 2 {
		t.Fatalf("Expected at least 2 FRQs in store, got %d", len(allRequests))
	}

	// Setup Feeder Topology
	nodeID := "transformer-001"
	topology := operator.FeederTopology{
		Capacities: []operator.NodeCapacity{
			{NodeID: nodeID, CapacityKW: 12.0}, // Limit to 12kW
		},
		Mappings: []operator.DeviceMapping{
			{DeviceLFDI: mrid1, NodeID: nodeID},
			{DeviceLFDI: mrid2, NodeID: nodeID},
		},
	}

	scheduler := operator.NewScheduler(allRequests)
	dispatcher := operator.NewFeederAwareDispatcher(scheduler, nil, topology)

	// Grid need: 15kW for 10 minutes
	gridReq := operator.GridServiceRequest{
		StartTime: now,
		Duration:  10 * time.Minute,
		PowerKW:   15.0,
	}

	scheduled, err := dispatcher.Schedule(gridReq)
	// We expect an error because 15kW > 12kW limit
	if err == nil {
		t.Errorf("Expected feeder limit error, got nil")
	}

	// Verify that the dispatcher correctly limited the dispatch
	totalPower := 0.0
	for _, ev := range scheduled {
		totalPower += float64(ev.Control.DERControlBase.OpModTargetW.Value) * 1e-3
	}

	if totalPower > 12.0 {
		t.Errorf("Feeder limit violated! Scheduled %.2f kW, limit 12.0 kW", totalPower)
	}

	t.Logf("Full integration flow successful. Scheduled %.2f kW within 12.0 kW feeder limit.", totalPower)
}
