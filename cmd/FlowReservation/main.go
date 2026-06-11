package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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
		Addr:              routes.FlowReservation,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig:         cfg,
		Handler:           tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	defer reg.Close()
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "FlowReservation.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /frq", http.HandlerFunc(h.GETFlowReservationRequestList))
	http.Handle("HEAD /frq", http.HandlerFunc(h.HEADFlowReservationRequestList))
	http.Handle("PUT /frq", http.HandlerFunc(h.PUTFlowReservationRequestList))
	http.Handle("POST /frq", http.HandlerFunc(h.POSTFlowReservationRequestList))
	http.Handle("DELETE /frq", http.HandlerFunc(h.DELETEFlowReservationRequestList))
	http.Handle("GET /frq/{id1}", http.HandlerFunc(h.GETFlowReservationRequest))
	http.Handle("HEAD /frq/{id1}", http.HandlerFunc(h.HEADFlowReservationRequest))
	http.Handle("PUT /frq/{id1}", http.HandlerFunc(h.PUTFlowReservationRequest))
	http.Handle("POST /frq/{id1}", http.HandlerFunc(h.POSTFlowReservationRequest))
	http.Handle("DELETE /frq/{id1}", http.HandlerFunc(h.DELETEFlowReservationRequest))
	http.Handle("GET /frp", http.HandlerFunc(h.GETFlowReservationResponseList))
	http.Handle("HEAD /frp", http.HandlerFunc(h.HEADFlowReservationResponseList))
	http.Handle("PUT /frp", http.HandlerFunc(h.PUTFlowReservationResponseList))
	http.Handle("POST /frp", http.HandlerFunc(h.POSTFlowReservationResponseList))
	http.Handle("DELETE /frp", http.HandlerFunc(h.DELETEFlowReservationResponseList))
	http.Handle("GET /frp/{id1}", http.HandlerFunc(h.GETFlowReservationResponse))
	http.Handle("HEAD /frp/{id1}", http.HandlerFunc(h.HEADFlowReservationResponse))
	http.Handle("PUT /frp/{id1}", http.HandlerFunc(h.PUTFlowReservationResponse))
	http.Handle("POST /frp/{id1}", http.HandlerFunc(h.POSTFlowReservationResponse))
	http.Handle("DELETE /frp/{id1}", http.HandlerFunc(h.DELETEFlowReservationResponse))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting FlowReservation on %s", routes.FlowReservation)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down FlowReservation server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
