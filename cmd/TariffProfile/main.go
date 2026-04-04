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
	http.HandleFunc("/", h.DELETEActiveTimeTariffIntervalList)
	http.HandleFunc("/", h.DELETEConsumptionTariffInterval)
	http.HandleFunc("/", h.DELETEConsumptionTariffIntervalList)
	http.HandleFunc("/", h.DELETERateComponent)
	http.HandleFunc("/", h.DELETERateComponentList)
	http.HandleFunc("/", h.DELETETariffProfile)
	http.HandleFunc("/", h.DELETETariffProfileList)
	http.HandleFunc("/", h.DELETETimeTariffInterval)
	http.HandleFunc("/", h.DELETETimeTariffIntervalList)
	http.HandleFunc("/", h.GETActiveTimeTariffIntervalList)
	http.HandleFunc("/", h.GETConsumptionTariffInterval)
	http.HandleFunc("/", h.GETConsumptionTariffIntervalList)
	http.HandleFunc("/", h.GETRateComponent)
	http.HandleFunc("/", h.GETRateComponentList)
	http.HandleFunc("/", h.GETTariffProfile)
	http.HandleFunc("/", h.GETTariffProfileList)
	http.HandleFunc("/", h.GETTimeTariffInterval)
	http.HandleFunc("/", h.GETTimeTariffIntervalList)
	http.HandleFunc("/", h.HEADActiveTimeTariffIntervalList)
	http.HandleFunc("/", h.HEADConsumptionTariffInterval)
	http.HandleFunc("/", h.HEADConsumptionTariffIntervalList)
	http.HandleFunc("/", h.HEADRateComponent)
	http.HandleFunc("/", h.HEADRateComponentList)
	http.HandleFunc("/", h.HEADTariffProfile)
	http.HandleFunc("/", h.HEADTariffProfileList)
	http.HandleFunc("/", h.HEADTimeTariffInterval)
	http.HandleFunc("/", h.HEADTimeTariffIntervalList)
	http.HandleFunc("/", h.POSTActiveTimeTariffIntervalList)
	http.HandleFunc("/", h.POSTConsumptionTariffInterval)
	http.HandleFunc("/", h.POSTConsumptionTariffIntervalList)
	http.HandleFunc("/", h.POSTRateComponent)
	http.HandleFunc("/", h.POSTRateComponentList)
	http.HandleFunc("/", h.POSTTariffProfile)
	http.HandleFunc("/", h.POSTTariffProfileList)
	http.HandleFunc("/", h.POSTTimeTariffInterval)
	http.HandleFunc("/", h.POSTTimeTariffIntervalList)
	http.HandleFunc("/", h.PUTActiveTimeTariffIntervalList)
	http.HandleFunc("/", h.PUTConsumptionTariffInterval)
	http.HandleFunc("/", h.PUTConsumptionTariffIntervalList)
	http.HandleFunc("/", h.PUTRateComponent)
	http.HandleFunc("/", h.PUTRateComponentList)
	http.HandleFunc("/", h.PUTTariffProfile)
	http.HandleFunc("/", h.PUTTariffProfileList)
	http.HandleFunc("/", h.PUTTimeTariffInterval)
	http.HandleFunc("/", h.PUTTimeTariffIntervalList)

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
