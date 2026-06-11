package main

import (
	"time"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/operator"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.Operator,
		TLSConfig: cfg,
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// In a real implementation, these would poll the FlowReservation and MUP services
	scheduler := operator.NewScheduler(nil)
	settlement := &operator.SettlementEngine{}

	http.HandleFunc("/operator/schedule", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var req operator.GridServiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		events, err := scheduler.Schedule(req)
		if err != nil {
			log.Printf("Scheduling warning: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})

	http.HandleFunc("/operator/settlement", func(w http.ResponseWriter, r *http.Request) {
		// Placeholder for settlement triggering
		report := settlement.CalculatePerformance("test-device", nil, nil)
		json.NewEncoder(w).Encode(report)
	})

	log.Printf("Starting operator on %s", routes.Operator)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
