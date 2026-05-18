package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/PPY/handler"
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
		Addr:      routes.PPY,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "PPY.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /ppy", http.HandlerFunc(h.GETPrepaymentList))
	http.Handle("HEAD /ppy", http.HandlerFunc(h.HEADPrepaymentList))
	http.Handle("PUT /ppy", http.HandlerFunc(h.PUTPrepaymentList))
	http.Handle("POST /ppy", http.HandlerFunc(h.POSTPrepaymentList))
	http.Handle("DELETE /ppy", http.HandlerFunc(h.DELETEPrepaymentList))
	http.Handle("GET /ppy/{id1}", http.HandlerFunc(h.GETPrepayment))
	http.Handle("HEAD /ppy/{id1}", http.HandlerFunc(h.HEADPrepayment))
	http.Handle("PUT /ppy/{id1}", http.HandlerFunc(h.PUTPrepayment))
	http.Handle("POST /ppy/{id1}", http.HandlerFunc(h.POSTPrepayment))
	http.Handle("DELETE /ppy/{id1}", http.HandlerFunc(h.DELETEPrepayment))
	http.Handle("GET /ppy/{id1}/ab", http.HandlerFunc(h.GETAccountBalance))
	http.Handle("HEAD /ppy/{id1}/ab", http.HandlerFunc(h.HEADAccountBalance))
	http.Handle("PUT /ppy/{id1}/ab", http.HandlerFunc(h.PUTAccountBalance))
	http.Handle("POST /ppy/{id1}/ab", http.HandlerFunc(h.POSTAccountBalance))
	http.Handle("DELETE /ppy/{id1}/ab", http.HandlerFunc(h.DELETEAccountBalance))
	http.Handle("GET /ppy/{id1}/os", http.HandlerFunc(h.GETPrepayOperationStatus))
	http.Handle("HEAD /ppy/{id1}/os", http.HandlerFunc(h.HEADPrepayOperationStatus))
	http.Handle("PUT /ppy/{id1}/os", http.HandlerFunc(h.PUTPrepayOperationStatus))
	http.Handle("POST /ppy/{id1}/os", http.HandlerFunc(h.POSTPrepayOperationStatus))
	http.Handle("DELETE /ppy/{id1}/os", http.HandlerFunc(h.DELETEPrepayOperationStatus))
	http.Handle("GET /ppy/{id1}/actsi", http.HandlerFunc(h.GETActiveSupplyInterruptionOverrideList))
	http.Handle("HEAD /ppy/{id1}/actsi", http.HandlerFunc(h.HEADActiveSupplyInterruptionOverrideList))
	http.Handle("PUT /ppy/{id1}/actsi", http.HandlerFunc(h.PUTActiveSupplyInterruptionOverrideList))
	http.Handle("POST /ppy/{id1}/actsi", http.HandlerFunc(h.POSTActiveSupplyInterruptionOverrideList))
	http.Handle("DELETE /ppy/{id1}/actsi", http.HandlerFunc(h.DELETEActiveSupplyInterruptionOverrideList))
	http.Handle("GET /ppy/{id1}/si", http.HandlerFunc(h.GETSupplyInterruptionOverrideList))
	http.Handle("HEAD /ppy/{id1}/si", http.HandlerFunc(h.HEADSupplyInterruptionOverrideList))
	http.Handle("PUT /ppy/{id1}/si", http.HandlerFunc(h.PUTSupplyInterruptionOverrideList))
	http.Handle("POST /ppy/{id1}/si", http.HandlerFunc(h.POSTSupplyInterruptionOverrideList))
	http.Handle("DELETE /ppy/{id1}/si", http.HandlerFunc(h.DELETESupplyInterruptionOverrideList))
	http.Handle("GET /ppy/{id1}/si/{id2}", http.HandlerFunc(h.GETSupplyInterruptionOverride))
	http.Handle("HEAD /ppy/{id1}/si/{id2}", http.HandlerFunc(h.HEADSupplyInterruptionOverride))
	http.Handle("PUT /ppy/{id1}/si/{id2}", http.HandlerFunc(h.PUTSupplyInterruptionOverride))
	http.Handle("POST /ppy/{id1}/si/{id2}", http.HandlerFunc(h.POSTSupplyInterruptionOverride))
	http.Handle("DELETE /ppy/{id1}/si/{id2}", http.HandlerFunc(h.DELETESupplyInterruptionOverride))
	http.Handle("GET /ppy/{id1}/cr", http.HandlerFunc(h.GETCreditRegisterList))
	http.Handle("HEAD /ppy/{id1}/cr", http.HandlerFunc(h.HEADCreditRegisterList))
	http.Handle("PUT /ppy/{id1}/cr", http.HandlerFunc(h.PUTCreditRegisterList))
	http.Handle("POST /ppy/{id1}/cr", http.HandlerFunc(h.POSTCreditRegisterList))
	http.Handle("DELETE /ppy/{id1}/cr", http.HandlerFunc(h.DELETECreditRegisterList))
	http.Handle("GET /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.GETCreditRegister))
	http.Handle("HEAD /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.HEADCreditRegister))
	http.Handle("PUT /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.PUTCreditRegister))
	http.Handle("POST /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.POSTCreditRegister))
	http.Handle("DELETE /ppy/{id1}/cr/{id2}", http.HandlerFunc(h.DELETECreditRegister))

	log.Printf("Starting PPY on %s", routes.PPY)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
