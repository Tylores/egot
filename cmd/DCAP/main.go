package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/DCAP/handler"
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
		Addr:      routes.DCAP,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "DCAP.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /dcap", http.HandlerFunc(h.GETDeviceCapability))
	http.Handle("HEAD /dcap", http.HandlerFunc(h.HEADDeviceCapability))
	http.Handle("PUT /dcap", http.HandlerFunc(h.PUTDeviceCapability))
	http.Handle("POST /dcap", http.HandlerFunc(h.POSTDeviceCapability))
	http.Handle("DELETE /dcap", http.HandlerFunc(h.DELETEDeviceCapability))

	log.Printf("Starting DCAP on %s", routes.DCAP)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
