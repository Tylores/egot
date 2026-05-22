package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/DER/handler"
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
		Addr:      routes.DER,
		TLSConfig: cfg,
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs

	repo := store.New(filepath.Join("data", "DER.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /derp", http.HandlerFunc(h.GETDERProgramList))
	http.Handle("HEAD /derp", http.HandlerFunc(h.HEADDERProgramList))
	http.Handle("PUT /derp", http.HandlerFunc(h.PUTDERProgramList))
	http.Handle("POST /derp", http.HandlerFunc(h.POSTDERProgramList))
	http.Handle("DELETE /derp", http.HandlerFunc(h.DELETEDERProgramList))
	http.Handle("GET /derp/{id1}", http.HandlerFunc(h.GETDERProgram))
	http.Handle("HEAD /derp/{id1}", http.HandlerFunc(h.HEADDERProgram))
	http.Handle("PUT /derp/{id1}", http.HandlerFunc(h.PUTDERProgram))
	http.Handle("POST /derp/{id1}", http.HandlerFunc(h.POSTDERProgram))
	http.Handle("DELETE /derp/{id1}", http.HandlerFunc(h.DELETEDERProgram))
	http.Handle("GET /derp/{id1}/actderc", http.HandlerFunc(h.GETActiveDERControlList))
	http.Handle("HEAD /derp/{id1}/actderc", http.HandlerFunc(h.HEADActiveDERControlList))
	http.Handle("PUT /derp/{id1}/actderc", http.HandlerFunc(h.PUTActiveDERControlList))
	http.Handle("POST /derp/{id1}/actderc", http.HandlerFunc(h.POSTActiveDERControlList))
	http.Handle("DELETE /derp/{id1}/actderc", http.HandlerFunc(h.DELETEActiveDERControlList))
	http.Handle("GET /derp/{id1}/derc", http.HandlerFunc(h.GETDERControlList))
	http.Handle("HEAD /derp/{id1}/derc", http.HandlerFunc(h.HEADDERControlList))
	http.Handle("PUT /derp/{id1}/derc", http.HandlerFunc(h.PUTDERControlList))
	http.Handle("POST /derp/{id1}/derc", http.HandlerFunc(h.POSTDERControlList))
	http.Handle("DELETE /derp/{id1}/derc", http.HandlerFunc(h.DELETEDERControlList))
	http.Handle("GET /derp/{id1}/derc/{id2}", http.HandlerFunc(h.GETDERControl))
	http.Handle("HEAD /derp/{id1}/derc/{id2}", http.HandlerFunc(h.HEADDERControl))
	http.Handle("PUT /derp/{id1}/derc/{id2}", http.HandlerFunc(h.PUTDERControl))
	http.Handle("POST /derp/{id1}/derc/{id2}", http.HandlerFunc(h.POSTDERControl))
	http.Handle("DELETE /derp/{id1}/derc/{id2}", http.HandlerFunc(h.DELETEDERControl))
	http.Handle("GET /derp/{id1}/dderc", http.HandlerFunc(h.GETDefaultDERControl))
	http.Handle("HEAD /derp/{id1}/dderc", http.HandlerFunc(h.HEADDefaultDERControl))
	http.Handle("PUT /derp/{id1}/dderc", http.HandlerFunc(h.PUTDefaultDERControl))
	http.Handle("POST /derp/{id1}/dderc", http.HandlerFunc(h.POSTDefaultDERControl))
	http.Handle("DELETE /derp/{id1}/dderc", http.HandlerFunc(h.DELETEDefaultDERControl))
	http.Handle("GET /derp/{id1}/dc", http.HandlerFunc(h.GETDERCurveList))
	http.Handle("HEAD /derp/{id1}/dc", http.HandlerFunc(h.HEADDERCurveList))
	http.Handle("PUT /derp/{id1}/dc", http.HandlerFunc(h.PUTDERCurveList))
	http.Handle("POST /derp/{id1}/dc", http.HandlerFunc(h.POSTDERCurveList))
	http.Handle("DELETE /derp/{id1}/dc", http.HandlerFunc(h.DELETEDERCurveList))
	http.Handle("GET /derp/{id1}/dc/{id2}", http.HandlerFunc(h.GETDERCurve))
	http.Handle("HEAD /derp/{id1}/dc/{id2}", http.HandlerFunc(h.HEADDERCurve))
	http.Handle("PUT /derp/{id1}/dc/{id2}", http.HandlerFunc(h.PUTDERCurve))
	http.Handle("POST /derp/{id1}/dc/{id2}", http.HandlerFunc(h.POSTDERCurve))
	http.Handle("DELETE /derp/{id1}/dc/{id2}", http.HandlerFunc(h.DELETEDERCurve))

	// EDevice DER sub-resource routes mapped to DER service
	http.Handle("GET /der", http.HandlerFunc(h.GETDERList))
	http.Handle("HEAD /der", http.HandlerFunc(h.HEADDERList))
	http.Handle("PUT /der", http.HandlerFunc(h.PUTDERList))
	http.Handle("POST /der", http.HandlerFunc(h.POSTDERList))
	http.Handle("DELETE /der", http.HandlerFunc(h.DELETEDERList))
	http.Handle("GET /der/{id1}", http.HandlerFunc(h.GETDER))
	http.Handle("HEAD /der/{id1}", http.HandlerFunc(h.HEADDER))
	http.Handle("PUT /der/{id1}", http.HandlerFunc(h.PUTDER))
	http.Handle("POST /der/{id1}", http.HandlerFunc(h.POSTDER))
	http.Handle("DELETE /der/{id1}", http.HandlerFunc(h.DELETEDER))
	http.Handle("GET /der/{id1}/upt", http.HandlerFunc(h.GETAssociatedUsagePoint))
	http.Handle("HEAD /der/{id1}/upt", http.HandlerFunc(h.HEADAssociatedUsagePoint))
	http.Handle("PUT /der/{id1}/upt", http.HandlerFunc(h.PUTAssociatedUsagePoint))
	http.Handle("POST /der/{id1}/upt", http.HandlerFunc(h.POSTAssociatedUsagePoint))
	http.Handle("DELETE /der/{id1}/upt", http.HandlerFunc(h.DELETEAssociatedUsagePoint))
	http.Handle("GET /der/{id1}/derp", http.HandlerFunc(h.GETAssociatedDERProgramList))
	http.Handle("HEAD /der/{id1}/derp", http.HandlerFunc(h.HEADAssociatedDERProgramList))
	http.Handle("PUT /der/{id1}/derp", http.HandlerFunc(h.PUTAssociatedDERProgramList))
	http.Handle("POST /der/{id1}/derp", http.HandlerFunc(h.POSTAssociatedDERProgramList))
	http.Handle("DELETE /der/{id1}/derp", http.HandlerFunc(h.DELETEAssociatedDERProgramList))
	http.Handle("GET /der/{id1}/cdc", http.HandlerFunc(h.GETCurrentDERControls))
	http.Handle("HEAD /der/{id1}/cdc", http.HandlerFunc(h.HEADCurrentDERControls))
	http.Handle("PUT /der/{id1}/cdc", http.HandlerFunc(h.PUTCurrentDERControls))
	http.Handle("POST /der/{id1}/cdc", http.HandlerFunc(h.POSTCurrentDERControls))
	http.Handle("DELETE /der/{id1}/cdc", http.HandlerFunc(h.DELETECurrentDERControls))
	http.Handle("GET /der/{id1}/cdp", http.HandlerFunc(h.GETCurrentDERProgram))
	http.Handle("HEAD /der/{id1}/cdp", http.HandlerFunc(h.HEADCurrentDERProgram))
	http.Handle("PUT /der/{id1}/cdp", http.HandlerFunc(h.PUTCurrentDERProgram))
	http.Handle("POST /der/{id1}/cdp", http.HandlerFunc(h.POSTCurrentDERProgram))
	http.Handle("DELETE /der/{id1}/cdp", http.HandlerFunc(h.DELETECurrentDERProgram))
	http.Handle("GET /der/{id1}/derg", http.HandlerFunc(h.GETDERSettings))
	http.Handle("HEAD /der/{id1}/derg", http.HandlerFunc(h.HEADDERSettings))
	http.Handle("PUT /der/{id1}/derg", http.HandlerFunc(h.PUTDERSettings))
	http.Handle("POST /der/{id1}/derg", http.HandlerFunc(h.POSTDERSettings))
	http.Handle("DELETE /der/{id1}/derg", http.HandlerFunc(h.DELETEDERSettings))
	http.Handle("GET /der/{id1}/ders", http.HandlerFunc(h.GETDERStatus))
	http.Handle("HEAD /der/{id1}/ders", http.HandlerFunc(h.HEADDERStatus))
	http.Handle("PUT /der/{id1}/ders", http.HandlerFunc(h.PUTDERStatus))
	http.Handle("POST /der/{id1}/ders", http.HandlerFunc(h.POSTDERStatus))
	http.Handle("DELETE /der/{id1}/ders", http.HandlerFunc(h.DELETEDERStatus))
	http.Handle("GET /der/{id1}/dera", http.HandlerFunc(h.GETDERAvailability))
	http.Handle("HEAD /der/{id1}/dera", http.HandlerFunc(h.HEADDERAvailability))
	http.Handle("PUT /der/{id1}/dera", http.HandlerFunc(h.PUTDERAvailability))
	http.Handle("POST /der/{id1}/dera", http.HandlerFunc(h.POSTDERAvailability))
	http.Handle("DELETE /der/{id1}/dera", http.HandlerFunc(h.DELETEDERAvailability))
	http.Handle("GET /der/{id1}/dercap", http.HandlerFunc(h.GETDERCapability))
	http.Handle("HEAD /der/{id1}/dercap", http.HandlerFunc(h.HEADDERCapability))
	http.Handle("PUT /der/{id1}/dercap", http.HandlerFunc(h.PUTDERCapability))
	http.Handle("POST /der/{id1}/dercap", http.HandlerFunc(h.POSTDERCapability))
	http.Handle("DELETE /der/{id1}/dercap", http.HandlerFunc(h.DELETEDERCapability))
	http.Handle("GET /der/{id1}/dercom", http.HandlerFunc(h.GETDERComponentList))
	http.Handle("HEAD /der/{id1}/dercom", http.HandlerFunc(h.HEADDERComponentList))
	http.Handle("PUT /der/{id1}/dercom", http.HandlerFunc(h.PUTDERComponentList))
	http.Handle("POST /der/{id1}/dercom", http.HandlerFunc(h.POSTDERComponentList))
	http.Handle("DELETE /der/{id1}/dercom", http.HandlerFunc(h.DELETEDERComponentList))
	http.Handle("GET /der/{id1}/dercom/{id2}", http.HandlerFunc(h.GETDERComponent))
	http.Handle("HEAD /der/{id1}/dercom/{id2}", http.HandlerFunc(h.HEADDERComponent))
	http.Handle("PUT /der/{id1}/dercom/{id2}", http.HandlerFunc(h.PUTDERComponent))
	http.Handle("POST /der/{id1}/dercom/{id2}", http.HandlerFunc(h.POSTDERComponent))
	http.Handle("DELETE /der/{id1}/dercom/{id2}", http.HandlerFunc(h.DELETEDERComponent))


	log.Printf("Starting DER on %s", routes.DER)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
