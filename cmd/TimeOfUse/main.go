package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/TimeOfUse/handler"
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
		Addr:      routes.TimeOfUse,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "TimeOfUse.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /tm", http.HandlerFunc(h.GETTime))
	http.Handle("HEAD /tm", http.HandlerFunc(h.HEADTime))
	http.Handle("PUT /tm", http.HandlerFunc(h.PUTTime))
	http.Handle("POST /tm", http.HandlerFunc(h.POSTTime))
	http.Handle("DELETE /tm", http.HandlerFunc(h.DELETETime))

	log.Printf("Starting TimeOfUse on %s", routes.TimeOfUse)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
