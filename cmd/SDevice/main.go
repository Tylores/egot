package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/SDevice/handler"
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
		Addr:      routes.SDevice,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs

	repo := store.New(filepath.Join("data", "SDevice.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /sdev", http.HandlerFunc(h.GETSelfDevice))
	http.Handle("HEAD /sdev", http.HandlerFunc(h.HEADSelfDevice))
	http.Handle("PUT /sdev", http.HandlerFunc(h.PUTSelfDevice))
	http.Handle("POST /sdev", http.HandlerFunc(h.POSTSelfDevice))
	http.Handle("DELETE /sdev", http.HandlerFunc(h.DELETESelfDevice))

	log.Printf("Starting SDevice on %s", routes.SDevice)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
