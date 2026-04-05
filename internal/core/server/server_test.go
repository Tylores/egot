package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/sep/uri"
)

func TestServerGet(t *testing.T) {
	// Run the HTTP server in a separate goroutine
	go func() {
		ServeHTTP(10)
	}()

	// Wait for server to be ready (up to 1 second)
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + routes.Core + uri.DeviceCapability)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	tests := []struct {
		method         string
		target         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			target:         uri.DeviceCapability,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			target:         uri.Time,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			target:         uri.EndDevice,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			target:         uri.EndDeviceList,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			target:         uri.Registration,
			expectedStatus: http.StatusUnauthorized,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		resp, err := http.Get("http://" + routes.Core + tt.target)

		if err != nil {
			t.Errorf("failed request: %v", err)
			continue
		}
		defer resp.Body.Close()

		// Assert the status code
		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				resp.StatusCode, tt.expectedStatus)
		}
	}
}
