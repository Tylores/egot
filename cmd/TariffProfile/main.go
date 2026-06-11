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

	"github.com/Tylores/egot/internal/TariffProfile/handler"
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
		Addr:              routes.TariffProfile,
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

	repo := store.New(filepath.Join("data", "TariffProfile.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /tp", http.HandlerFunc(h.GETTariffProfileList))
	http.Handle("HEAD /tp", http.HandlerFunc(h.HEADTariffProfileList))
	http.Handle("PUT /tp", http.HandlerFunc(h.PUTTariffProfileList))
	http.Handle("POST /tp", http.HandlerFunc(h.POSTTariffProfileList))
	http.Handle("DELETE /tp", http.HandlerFunc(h.DELETETariffProfileList))
	http.Handle("GET /tp/{id1}", http.HandlerFunc(h.GETTariffProfile))
	http.Handle("HEAD /tp/{id1}", http.HandlerFunc(h.HEADTariffProfile))
	http.Handle("PUT /tp/{id1}", http.HandlerFunc(h.PUTTariffProfile))
	http.Handle("POST /tp/{id1}", http.HandlerFunc(h.POSTTariffProfile))
	http.Handle("DELETE /tp/{id1}", http.HandlerFunc(h.DELETETariffProfile))
	http.Handle("GET /tp/{id1}/rc", http.HandlerFunc(h.GETRateComponentList))
	http.Handle("HEAD /tp/{id1}/rc", http.HandlerFunc(h.HEADRateComponentList))
	http.Handle("PUT /tp/{id1}/rc", http.HandlerFunc(h.PUTRateComponentList))
	http.Handle("POST /tp/{id1}/rc", http.HandlerFunc(h.POSTRateComponentList))
	http.Handle("DELETE /tp/{id1}/rc", http.HandlerFunc(h.DELETERateComponentList))
	http.Handle("GET /tp/{id1}/rc/{id2}", http.HandlerFunc(h.GETRateComponent))
	http.Handle("HEAD /tp/{id1}/rc/{id2}", http.HandlerFunc(h.HEADRateComponent))
	http.Handle("PUT /tp/{id1}/rc/{id2}", http.HandlerFunc(h.PUTRateComponent))
	http.Handle("POST /tp/{id1}/rc/{id2}", http.HandlerFunc(h.POSTRateComponent))
	http.Handle("DELETE /tp/{id1}/rc/{id2}", http.HandlerFunc(h.DELETERateComponent))
	http.Handle("GET /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.GETActiveTimeTariffIntervalList))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.HEADActiveTimeTariffIntervalList))
	http.Handle("PUT /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.PUTActiveTimeTariffIntervalList))
	http.Handle("POST /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.POSTActiveTimeTariffIntervalList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.DELETEActiveTimeTariffIntervalList))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.GETTimeTariffIntervalList))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.HEADTimeTariffIntervalList))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.PUTTimeTariffIntervalList))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.POSTTimeTariffIntervalList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.DELETETimeTariffIntervalList))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.GETTimeTariffInterval))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.HEADTimeTariffInterval))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.PUTTimeTariffInterval))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.POSTTimeTariffInterval))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.DELETETimeTariffInterval))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.GETConsumptionTariffIntervalList))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.HEADConsumptionTariffIntervalList))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.PUTConsumptionTariffIntervalList))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.POSTConsumptionTariffIntervalList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.DELETEConsumptionTariffIntervalList))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.GETConsumptionTariffInterval))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.HEADConsumptionTariffInterval))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.PUTConsumptionTariffInterval))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.POSTConsumptionTariffInterval))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.DELETEConsumptionTariffInterval))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting TariffProfile on %s", routes.TariffProfile)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down TariffProfile server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
