package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/BRS/handler"
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
		Addr:      routes.BRS,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "BRS.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /brs", http.HandlerFunc(h.GETBillingReadingSetList))
	http.Handle("HEAD /brs", http.HandlerFunc(h.HEADBillingReadingSetList))
	http.Handle("PUT /brs", http.HandlerFunc(h.PUTBillingReadingSetList))
	http.Handle("POST /brs", http.HandlerFunc(h.POSTBillingReadingSetList))
	http.Handle("DELETE /brs", http.HandlerFunc(h.DELETEBillingReadingSetList))
	http.Handle("GET /brs/{id1}", http.HandlerFunc(h.GETBillingReadingSet))
	http.Handle("HEAD /brs/{id1}", http.HandlerFunc(h.HEADBillingReadingSet))
	http.Handle("PUT /brs/{id1}", http.HandlerFunc(h.PUTBillingReadingSet))
	http.Handle("POST /brs/{id1}", http.HandlerFunc(h.POSTBillingReadingSet))
	http.Handle("DELETE /brs/{id1}", http.HandlerFunc(h.DELETEBillingReadingSet))
	http.Handle("GET /brs/{id1}/br", http.HandlerFunc(h.GETBillingReadingList))
	http.Handle("HEAD /brs/{id1}/br", http.HandlerFunc(h.HEADBillingReadingList))
	http.Handle("PUT /brs/{id1}/br", http.HandlerFunc(h.PUTBillingReadingList))
	http.Handle("POST /brs/{id1}/br", http.HandlerFunc(h.POSTBillingReadingList))
	http.Handle("DELETE /brs/{id1}/br", http.HandlerFunc(h.DELETEBillingReadingList))
	http.Handle("GET /brs/{id1}/br/{id2}", http.HandlerFunc(h.GETBillingReading))
	http.Handle("HEAD /brs/{id1}/br/{id2}", http.HandlerFunc(h.HEADBillingReading))
	http.Handle("PUT /brs/{id1}/br/{id2}", http.HandlerFunc(h.PUTBillingReading))
	http.Handle("POST /brs/{id1}/br/{id2}", http.HandlerFunc(h.POSTBillingReading))
	http.Handle("DELETE /brs/{id1}/br/{id2}", http.HandlerFunc(h.DELETEBillingReading))

	log.Printf("Starting BRS on %s", routes.BRS)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
