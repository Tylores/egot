package server

import (
	"net/http"
	"testing"

	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/sep/uri"
)

func TestServerGet(t *testing.T) {
	// Run the HTTP server in a separate goroutine
	go func() {
		ServeHTTP(10)
	}()

	tests := []struct {
		method         string
		target         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			target:         uri.DeviceCapability,
			expectedStatus: http.StatusOK,
		},
		{
			target:         uri.Time,
			expectedStatus: http.StatusOK,
		},
		{
			target:         uri.EndDevice,
			expectedStatus: http.StatusOK,
		},
		{
			target:         uri.EndDeviceList,
			expectedStatus: http.StatusOK,
		},
		{
			target:         uri.Registration,
			expectedStatus: http.StatusOK,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		resp, err := http.Get("http://" + routes.Core + tt.target)

		if err != nil {
			t.Errorf("failed request: %v", err)
		}

		// Assert the status code
		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				resp.StatusCode, tt.expectedStatus)
		}
	}
}
