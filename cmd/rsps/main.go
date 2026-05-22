package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/rsps/handler"
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
		Addr:      routes.Rsps,
		TLSConfig: cfg,
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs

	repo := store.New(filepath.Join("data", "rsps.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /rsps", http.HandlerFunc(h.GETResponseSetList))
	http.Handle("HEAD /rsps", http.HandlerFunc(h.HEADResponseSetList))
	http.Handle("PUT /rsps", http.HandlerFunc(h.PUTResponseSetList))
	http.Handle("POST /rsps", http.HandlerFunc(h.POSTResponseSetList))
	http.Handle("DELETE /rsps", http.HandlerFunc(h.DELETEResponseSetList))
	http.Handle("GET /rsps/{id1}", http.HandlerFunc(h.GETResponseSet))
	http.Handle("HEAD /rsps/{id1}", http.HandlerFunc(h.HEADResponseSet))
	http.Handle("PUT /rsps/{id1}", http.HandlerFunc(h.PUTResponseSet))
	http.Handle("POST /rsps/{id1}", http.HandlerFunc(h.POSTResponseSet))
	http.Handle("DELETE /rsps/{id1}", http.HandlerFunc(h.DELETEResponseSet))
	http.Handle("GET /rsps/{id1}/rsp", http.HandlerFunc(h.GETResponseList))
	http.Handle("HEAD /rsps/{id1}/rsp", http.HandlerFunc(h.HEADResponseList))
	http.Handle("PUT /rsps/{id1}/rsp", http.HandlerFunc(h.PUTResponseList))
	http.Handle("POST /rsps/{id1}/rsp", http.HandlerFunc(h.POSTResponseList))
	http.Handle("DELETE /rsps/{id1}/rsp", http.HandlerFunc(h.DELETEResponseList))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEResponse))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETPriceResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADPriceResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTPriceResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTPriceResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEPriceResponse))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETTextResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADTextResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTTextResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTTextResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETETextResponse))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETDefaultDERControlResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADDefaultDERControlResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTDefaultDERControlResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTDefaultDERControlResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEDefaultDERControlResponse))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETDERControlResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADDERControlResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTDERControlResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTDERControlResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEDERControlResponse))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETFlowReservationResponseResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADFlowReservationResponseResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTFlowReservationResponseResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTFlowReservationResponseResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEFlowReservationResponseResponse))
	http.Handle("GET /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.GETDrResponse))
	http.Handle("HEAD /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.HEADDrResponse))
	http.Handle("PUT /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.PUTDrResponse))
	http.Handle("POST /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.POSTDrResponse))
	http.Handle("DELETE /rsps/{id1}/rsp/{id2}", http.HandlerFunc(h.DELETEDrResponse))

	log.Printf("Starting rsps on %s", routes.Rsps)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
