package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/TariffProfile/handler"
"github.com/Tylores/egot/internal/TariffProfile/repository/memory"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      "egot.internal.com:8022",
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)

// Register handler methods
	http.Handle("DELETE /tp", http.HandlerFunc(h.DELETETariffProfileList))
	http.Handle("GET /tp", http.HandlerFunc(h.GETTariffProfileList))
	http.Handle("HEAD /tp", http.HandlerFunc(h.HEADTariffProfileList))
	http.Handle("POST /tp", http.HandlerFunc(h.POSTTariffProfileList))
	http.Handle("PUT /tp", http.HandlerFunc(h.PUTTariffProfileList))
	http.Handle("DELETE /tp/{id1}", http.HandlerFunc(h.DELETETariffProfile))
	http.Handle("GET /tp/{id1}", http.HandlerFunc(h.GETTariffProfile))
	http.Handle("HEAD /tp/{id1}", http.HandlerFunc(h.HEADTariffProfile))
	http.Handle("POST /tp/{id1}", http.HandlerFunc(h.POSTTariffProfile))
	http.Handle("PUT /tp/{id1}", http.HandlerFunc(h.PUTTariffProfile))
	http.Handle("DELETE /tp/{id1}/rc", http.HandlerFunc(h.DELETERateComponentList))
	http.Handle("GET /tp/{id1}/rc", http.HandlerFunc(h.GETRateComponentList))
	http.Handle("HEAD /tp/{id1}/rc", http.HandlerFunc(h.HEADRateComponentList))
	http.Handle("POST /tp/{id1}/rc", http.HandlerFunc(h.POSTRateComponentList))
	http.Handle("PUT /tp/{id1}/rc", http.HandlerFunc(h.PUTRateComponentList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}", http.HandlerFunc(h.DELETERateComponent))
	http.Handle("GET /tp/{id1}/rc/{id2}", http.HandlerFunc(h.GETRateComponent))
	http.Handle("HEAD /tp/{id1}/rc/{id2}", http.HandlerFunc(h.HEADRateComponent))
	http.Handle("POST /tp/{id1}/rc/{id2}", http.HandlerFunc(h.POSTRateComponent))
	http.Handle("PUT /tp/{id1}/rc/{id2}", http.HandlerFunc(h.PUTRateComponent))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.DELETEActiveTimeTariffIntervalList))
	http.Handle("GET /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.GETActiveTimeTariffIntervalList))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.HEADActiveTimeTariffIntervalList))
	http.Handle("POST /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.POSTActiveTimeTariffIntervalList))
	http.Handle("PUT /tp/{id1}/rc/{id2}/acttti", http.HandlerFunc(h.PUTActiveTimeTariffIntervalList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.DELETETimeTariffIntervalList))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.GETTimeTariffIntervalList))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.HEADTimeTariffIntervalList))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.POSTTimeTariffIntervalList))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti", http.HandlerFunc(h.PUTTimeTariffIntervalList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.DELETETimeTariffInterval))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.GETTimeTariffInterval))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.HEADTimeTariffInterval))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.POSTTimeTariffInterval))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}", http.HandlerFunc(h.PUTTimeTariffInterval))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.DELETEConsumptionTariffIntervalList))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.GETConsumptionTariffIntervalList))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.HEADConsumptionTariffIntervalList))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.POSTConsumptionTariffIntervalList))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}/cti", http.HandlerFunc(h.PUTConsumptionTariffIntervalList))
	http.Handle("DELETE /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.DELETEConsumptionTariffInterval))
	http.Handle("GET /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.GETConsumptionTariffInterval))
	http.Handle("HEAD /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.HEADConsumptionTariffInterval))
	http.Handle("POST /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.POSTConsumptionTariffInterval))
	http.Handle("PUT /tp/{id1}/rc/{id2}/tti/{id3}/cti/{id4}", http.HandlerFunc(h.PUTConsumptionTariffInterval))
err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
