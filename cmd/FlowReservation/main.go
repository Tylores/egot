package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/FlowReservation/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.FlowReservation,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "FlowReservation.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /edev/{id1}/frq", http.HandlerFunc(h.GETFlowReservationRequestList))
	http.Handle("HEAD /edev/{id1}/frq", http.HandlerFunc(h.HEADFlowReservationRequestList))
	http.Handle("PUT /edev/{id1}/frq", http.HandlerFunc(h.PUTFlowReservationRequestList))
	http.Handle("POST /edev/{id1}/frq", http.HandlerFunc(h.POSTFlowReservationRequestList))
	http.Handle("DELETE /edev/{id1}/frq", http.HandlerFunc(h.DELETEFlowReservationRequestList))
	http.Handle("GET /edev/{id1}/frq/{id2}", http.HandlerFunc(h.GETFlowReservationRequest))
	http.Handle("HEAD /edev/{id1}/frq/{id2}", http.HandlerFunc(h.HEADFlowReservationRequest))
	http.Handle("PUT /edev/{id1}/frq/{id2}", http.HandlerFunc(h.PUTFlowReservationRequest))
	http.Handle("POST /edev/{id1}/frq/{id2}", http.HandlerFunc(h.POSTFlowReservationRequest))
	http.Handle("DELETE /edev/{id1}/frq/{id2}", http.HandlerFunc(h.DELETEFlowReservationRequest))

	log.Printf("Starting FlowReservation on %s", routes.FlowReservation)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
