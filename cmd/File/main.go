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

	"github.com/Tylores/egot/internal/File/handler"
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
		Addr:      routes.File,
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

	repo := store.New(filepath.Join("data", "File.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /file", http.HandlerFunc(h.GETFileList))
	http.Handle("HEAD /file", http.HandlerFunc(h.HEADFileList))
	http.Handle("PUT /file", http.HandlerFunc(h.PUTFileList))
	http.Handle("POST /file", http.HandlerFunc(h.POSTFileList))
	http.Handle("DELETE /file", http.HandlerFunc(h.DELETEFileList))
	http.Handle("GET /file/{id1}", http.HandlerFunc(h.GETFile))
	http.Handle("HEAD /file/{id1}", http.HandlerFunc(h.HEADFile))
	http.Handle("PUT /file/{id1}", http.HandlerFunc(h.PUTFile))
	http.Handle("POST /file/{id1}", http.HandlerFunc(h.POSTFile))
	http.Handle("DELETE /file/{id1}", http.HandlerFunc(h.DELETEFile))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting File on %s", routes.File)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down File server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
