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

	"github.com/Tylores/egot/internal/DR/handler"
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
		Addr:              routes.DR,
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

	repo := store.New(filepath.Join("data", "DR.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /dr", http.HandlerFunc(h.GETDemandResponseProgramList))
	http.Handle("HEAD /dr", http.HandlerFunc(h.HEADDemandResponseProgramList))
	http.Handle("PUT /dr", http.HandlerFunc(h.PUTDemandResponseProgramList))
	http.Handle("POST /dr", http.HandlerFunc(h.POSTDemandResponseProgramList))
	http.Handle("DELETE /dr", http.HandlerFunc(h.DELETEDemandResponseProgramList))
	http.Handle("GET /dr/{id1}", http.HandlerFunc(h.GETDemandResponseProgram))
	http.Handle("HEAD /dr/{id1}", http.HandlerFunc(h.HEADDemandResponseProgram))
	http.Handle("PUT /dr/{id1}", http.HandlerFunc(h.PUTDemandResponseProgram))
	http.Handle("POST /dr/{id1}", http.HandlerFunc(h.POSTDemandResponseProgram))
	http.Handle("DELETE /dr/{id1}", http.HandlerFunc(h.DELETEDemandResponseProgram))
	http.Handle("GET /dr/{id1}/actedc", http.HandlerFunc(h.GETActiveEndDeviceControlList))
	http.Handle("HEAD /dr/{id1}/actedc", http.HandlerFunc(h.HEADActiveEndDeviceControlList))
	http.Handle("PUT /dr/{id1}/actedc", http.HandlerFunc(h.PUTActiveEndDeviceControlList))
	http.Handle("POST /dr/{id1}/actedc", http.HandlerFunc(h.POSTActiveEndDeviceControlList))
	http.Handle("DELETE /dr/{id1}/actedc", http.HandlerFunc(h.DELETEActiveEndDeviceControlList))
	http.Handle("GET /dr/{id1}/edc", http.HandlerFunc(h.GETEndDeviceControlList))
	http.Handle("HEAD /dr/{id1}/edc", http.HandlerFunc(h.HEADEndDeviceControlList))
	http.Handle("PUT /dr/{id1}/edc", http.HandlerFunc(h.PUTEndDeviceControlList))
	http.Handle("POST /dr/{id1}/edc", http.HandlerFunc(h.POSTEndDeviceControlList))
	http.Handle("DELETE /dr/{id1}/edc", http.HandlerFunc(h.DELETEEndDeviceControlList))
	http.Handle("GET /dr/{id1}/edc/{id2}", http.HandlerFunc(h.GETEndDeviceControl))
	http.Handle("HEAD /dr/{id1}/edc/{id2}", http.HandlerFunc(h.HEADEndDeviceControl))
	http.Handle("PUT /dr/{id1}/edc/{id2}", http.HandlerFunc(h.PUTEndDeviceControl))
	http.Handle("POST /dr/{id1}/edc/{id2}", http.HandlerFunc(h.POSTEndDeviceControl))
	http.Handle("DELETE /dr/{id1}/edc/{id2}", http.HandlerFunc(h.DELETEEndDeviceControl))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting DR on %s", routes.DR)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down DR server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
