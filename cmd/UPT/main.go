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

	"github.com/Tylores/egot/internal/UPT/handler"
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
		Addr:      routes.UPT,
		TLSConfig: cfg,
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	defer reg.Close()
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "UPT.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /upt", http.HandlerFunc(h.GETUsagePointList))
	http.Handle("HEAD /upt", http.HandlerFunc(h.HEADUsagePointList))
	http.Handle("PUT /upt", http.HandlerFunc(h.PUTUsagePointList))
	http.Handle("POST /upt", http.HandlerFunc(h.POSTUsagePointList))
	http.Handle("DELETE /upt", http.HandlerFunc(h.DELETEUsagePointList))
	http.Handle("GET /upt/{id1}", http.HandlerFunc(h.GETUsagePoint))
	http.Handle("HEAD /upt/{id1}", http.HandlerFunc(h.HEADUsagePoint))
	http.Handle("PUT /upt/{id1}", http.HandlerFunc(h.PUTUsagePoint))
	http.Handle("POST /upt/{id1}", http.HandlerFunc(h.POSTUsagePoint))
	http.Handle("DELETE /upt/{id1}", http.HandlerFunc(h.DELETEUsagePoint))
	http.Handle("GET /upt/{id1}/mr", http.HandlerFunc(h.GETMeterReadingList))
	http.Handle("HEAD /upt/{id1}/mr", http.HandlerFunc(h.HEADMeterReadingList))
	http.Handle("PUT /upt/{id1}/mr", http.HandlerFunc(h.PUTMeterReadingList))
	http.Handle("POST /upt/{id1}/mr", http.HandlerFunc(h.POSTMeterReadingList))
	http.Handle("DELETE /upt/{id1}/mr", http.HandlerFunc(h.DELETEMeterReadingList))
	http.Handle("GET /upt/{id1}/mr/{id2}", http.HandlerFunc(h.GETMeterReading))
	http.Handle("HEAD /upt/{id1}/mr/{id2}", http.HandlerFunc(h.HEADMeterReading))
	http.Handle("PUT /upt/{id1}/mr/{id2}", http.HandlerFunc(h.PUTMeterReading))
	http.Handle("POST /upt/{id1}/mr/{id2}", http.HandlerFunc(h.POSTMeterReading))
	http.Handle("DELETE /upt/{id1}/mr/{id2}", http.HandlerFunc(h.DELETEMeterReading))
	http.Handle("GET /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.GETReadingType))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.HEADReadingType))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.PUTReadingType))
	http.Handle("POST /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.POSTReadingType))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.DELETEReadingType))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.GETReadingSetList))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.HEADReadingSetList))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.PUTReadingSetList))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.POSTReadingSetList))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.DELETEReadingSetList))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.GETReadingSet))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.HEADReadingSet))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.PUTReadingSet))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.POSTReadingSet))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.DELETEReadingSet))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.GETReadingList))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.HEADReadingList))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.PUTReadingList))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.POSTReadingList))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.DELETEReadingList))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.GETReading))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.HEADReading))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.PUTReading))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.POSTReading))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.DELETEReading))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting UPT on %s", routes.UPT)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down UPT server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
