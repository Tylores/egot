package main

import (
	"log"
	"net/http"
	"path/filepath"

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
		Addr:      routes.MUP,
		TLSConfig: cfg,
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs

	repo := store.New(filepath.Join("data", "MUP.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

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

	log.Printf("Starting MUP on %s", routes.MUP)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
