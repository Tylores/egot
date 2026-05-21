package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/DR/handler"
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
		Addr:      routes.DR,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs

	repo := store.New(filepath.Join("data", "DR.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /dr", http.HandlerFunc(h.GETDemandResponseProgramList))
	http.Handle("HEAD /dr", http.HandlerFunc(h.HEADDemandResponseProgramList))
	http.Handle("PUT /dr", http.HandlerFunc(h.PUTDemandResponseProgramList))
	http.Handle("POST /dr", http.HandlerFunc(h.POSTDemandResponseProgramList))
	http.Handle("DELETE /dr", http.HandlerFunc(h.DELETEDemandResponseProgramList))
	http.Handle("GET /dr/{id1}", http.HandlerFunc(h.GETDemandResponseProgram))
	http.Handle("HEAD /dr/{id1}", http.HandlerFunc(h.HEADDemandResponseProgram))
	http.Handle("PUT /dr/{id1}", http.HandlerFunc(h.PUTDemandResponseProgram))
	http.Handle("POST /dr/{id1}", http.HandlerFunc(h.POSTDemandResponseProgram))
	http.Handle("DELETE /dr/{id1}", http.HandlerFunc(h.DELETEDemandResponseProgram))
	http.Handle("GET /dr/{id1}/actedc", http.HandlerFunc(h.GETActiveEndDeviceControlList))
	http.Handle("HEAD /dr/{id1}/actedc", http.HandlerFunc(h.HEADActiveEndDeviceControlList))
	http.Handle("PUT /dr/{id1}/actedc", http.HandlerFunc(h.PUTActiveEndDeviceControlList))
	http.Handle("POST /dr/{id1}/actedc", http.HandlerFunc(h.POSTActiveEndDeviceControlList))
	http.Handle("DELETE /dr/{id1}/actedc", http.HandlerFunc(h.DELETEActiveEndDeviceControlList))
	http.Handle("GET /dr/{id1}/edc", http.HandlerFunc(h.GETEndDeviceControlList))
	http.Handle("HEAD /dr/{id1}/edc", http.HandlerFunc(h.HEADEndDeviceControlList))
	http.Handle("PUT /dr/{id1}/edc", http.HandlerFunc(h.PUTEndDeviceControlList))
	http.Handle("POST /dr/{id1}/edc", http.HandlerFunc(h.POSTEndDeviceControlList))
	http.Handle("DELETE /dr/{id1}/edc", http.HandlerFunc(h.DELETEEndDeviceControlList))
	http.Handle("GET /dr/{id1}/edc/{id2}", http.HandlerFunc(h.GETEndDeviceControl))
	http.Handle("HEAD /dr/{id1}/edc/{id2}", http.HandlerFunc(h.HEADEndDeviceControl))
	http.Handle("PUT /dr/{id1}/edc/{id2}", http.HandlerFunc(h.PUTEndDeviceControl))
	http.Handle("POST /dr/{id1}/edc/{id2}", http.HandlerFunc(h.POSTEndDeviceControl))
	http.Handle("DELETE /dr/{id1}/edc/{id2}", http.HandlerFunc(h.DELETEEndDeviceControl))

	log.Printf("Starting DR on %s", routes.DR)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
