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

	"github.com/Tylores/egot/internal/MUP/handler"
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
		Addr:              routes.MUP,
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

	repo := store.New(filepath.Join("data", "MUP.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /mup", http.HandlerFunc(h.GETMirrorUsagePointList))
	http.Handle("HEAD /mup", http.HandlerFunc(h.HEADMirrorUsagePointList))
	http.Handle("PUT /mup", http.HandlerFunc(h.PUTMirrorUsagePointList))
	http.Handle("POST /mup", http.HandlerFunc(h.POSTMirrorUsagePointList))
	http.Handle("DELETE /mup", http.HandlerFunc(h.DELETEMirrorUsagePointList))
	http.Handle("GET /mup/{id1}", http.HandlerFunc(h.GETMirrorUsagePoint))
	http.Handle("HEAD /mup/{id1}", http.HandlerFunc(h.HEADMirrorUsagePoint))
	http.Handle("PUT /mup/{id1}", http.HandlerFunc(h.PUTMirrorUsagePoint))
	http.Handle("POST /mup/{id1}", http.HandlerFunc(h.POSTMirrorUsagePoint))
	http.Handle("DELETE /mup/{id1}", http.HandlerFunc(h.DELETEMirrorUsagePoint))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting MUP on %s", routes.MUP)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down MUP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
