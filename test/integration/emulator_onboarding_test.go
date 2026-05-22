package integration

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Tylores/egot/sep"
)

func TestEmulatorOnboardingAndCSIPStatus(t *testing.T) {
	var mu sync.Mutex
	receivedPaths := make(map[string][]string)
	postedStatuses := make(map[uint8]bool)

	// Create test server with all required IEEE 2030.5 endpoints
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		receivedPaths[req.URL.Path] = append(receivedPaths[req.URL.Path], req.Method)

		switch {
		case req.URL.Path == "/dcap":
			w.Header().Set("Content-Type", sep.ContentType)
			xml.NewEncoder(w).Encode(&sep.DeviceCapability{
				EndDeviceListLink: &sep.EndDeviceListLink{
					ListLink: &sep.ListLink{
						Link: &sep.Link{
							HrefAttr: "/edev",
						},
					},
				},
			})

		case req.URL.Path == "/tm":
			w.Header().Set("Content-Type", sep.ContentType)
			xml.NewEncoder(w).Encode(&sep.Time{
				CurrentTime: sep.TimeType(time.Now().Unix()),
				LocalTime:   sep.TimeType(time.Now().Unix()),
				Quality:     7,
			})

		case req.URL.Path == "/edev":
			w.Header().Set("Location", "/edev/12345")
			w.WriteHeader(http.StatusCreated)

		case req.URL.Path == "/edev/12345/rg":
			w.Header().Set("Content-Type", sep.ContentType)
			pin := sep.PINType(111115)
			xml.NewEncoder(w).Encode(&sep.Registration{
				PIN: &pin,
			})

		case req.URL.Path == "/mup":
			w.Header().Set("Location", "/mup/1")
			w.WriteHeader(http.StatusCreated)

		case req.URL.Path == "/mup/1":
			w.WriteHeader(http.StatusNoContent)

		case req.URL.Path == "/derp/1/actderc":
			w.Header().Set("Content-Type", sep.ContentType)
			now := time.Now().Unix()
			nowVal := sep.TimeType(now)
			mrid := "mock-control-001"
			xml.NewEncoder(w).Encode(&sep.DERControlList{
				SubscribableList: &sep.SubscribableList{
					SubscribableResource: &sep.SubscribableResource{
						Resource: &sep.Resource{
							HrefAttr: "/derp/1/actderc",
						},
					},
					AllAttr:     1,
					ResultsAttr: 1,
				},
				DERControl: []*sep.DERControl{
					{
						RandomizableEvent: &sep.RandomizableEvent{
							Event: &sep.Event{
								RespondableSubscribableIdentifiedObject: &sep.RespondableSubscribableIdentifiedObject{
									RespondableResource: &sep.RespondableResource{
										Resource: &sep.Resource{
											HrefAttr: "/derp/1/derc/1",
										},
										ResponseRequiredAttr: "03",
									},
									MRID: &sep.MRIDType{HexBinary128: &mrid},
								},
								Interval: &sep.DateTimeInterval{
									Start:    &nowVal,
									Duration: 2, // 2 seconds duration
								},
							},
						},
						DERControlBase: &sep.DERControlBase{
							OpModFixedW: &sep.SignedPerCentControlType{
								SignedPerCent: func() *sep.SignedPerCent { v := sep.SignedPerCent(5000); return &v }(),
							},
						},
					},
				},
			})

		case strings.HasPrefix(req.URL.Path, "/rsps/") && strings.HasSuffix(req.URL.Path, "/rsp"):
			var body sep.DERControlResponse
			dec := xml.NewDecoder(req.Body)
			if err := dec.Decode(&body); err == nil && body.Response != nil {
				postedStatuses[body.Response.Status] = true
			}
			w.WriteHeader(http.StatusCreated)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Parse host/port from the test TLS server
	serverAddr := server.Listener.Addr().String()

	// Launch emulator-der as a subprocess
	cmd := exec.Command("./bin/emulator-der",
		"-type", "ess",
		"-name", "test-emu",
		"-gateway", serverAddr,
		"-ssl", "./ssl",
		"-interval", "1s",
		"-count", "1",
	)
	cmd.Dir = "/home/tylor/phd/egot"

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start emulator subprocess: %v", err)
	}

	// Let the emulator run for 4 seconds so that:
	// - Onboarding succeeds
	// - Received and Started statuses are posted
	// - Control completes (2s duration) and Completed status is posted
	time.Sleep(4 * time.Second)

	// Terminate the emulator process
	_ = cmd.Process.Kill()
	_ = cmd.Wait()

	outputLog := outBuf.String()
	t.Logf("Emulator logs:\n%s", outputLog)

	// Verify onboarding steps in log output
	expectedLogs := []string{
		"[CSIP] Step 1: Discovery GET",
		"[CSIP] Step 2: Time Sync GET",
		"[CSIP] Step 3: Device Registration POST",
		"[CSIP] Step 4: PIN Verification GET",
		"CSIP Onboarding completed successfully",
		"posted control response status 1",
		"posted control response status 2",
		"posted control response status 3",
	}

	for _, logMsg := range expectedLogs {
		if !strings.Contains(outputLog, logMsg) {
			t.Errorf("Expected emulator logs to contain %q, but it was missing", logMsg)
		}
	}

	// Verify server received requests
	mu.Lock()
	defer mu.Unlock()

	requiredPaths := []string{"/dcap", "/tm", "/edev", "/mup"}
	for _, p := range requiredPaths {
		if _, exists := receivedPaths[p]; !exists {
			t.Errorf("Expected mock server to receive request on path %q, but none received", p)
		}
	}

	// Verify posted statuses (Received = 1, Started = 2, Completed = 3)
	for _, status := range []uint8{1, 2, 3} {
		if !postedStatuses[status] {
			t.Errorf("Expected mock server to receive posted control status %d, but it was missing", status)
		}
	}
}
