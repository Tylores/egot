package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/UPT/handler"
"github.com/Tylores/egot/internal/UPT/repository/memory"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      "egot.internal.com:8024",
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)

// Register handler methods
	http.HandleFunc("/", h.DELETEMeterReading)
	http.HandleFunc("/", h.DELETEMeterReadingList)
	http.HandleFunc("/", h.DELETEReading)
	http.HandleFunc("/", h.DELETEReadingList)
	http.HandleFunc("/", h.DELETEReadingSet)
	http.HandleFunc("/", h.DELETEReadingSetList)
	http.HandleFunc("/", h.DELETEReadingType)
	http.HandleFunc("/", h.DELETEUsagePoint)
	http.HandleFunc("/", h.DELETEUsagePointList)
	http.HandleFunc("/", h.GETMeterReading)
	http.HandleFunc("/", h.GETMeterReadingList)
	http.HandleFunc("/", h.GETReading)
	http.HandleFunc("/", h.GETReadingList)
	http.HandleFunc("/", h.GETReadingSet)
	http.HandleFunc("/", h.GETReadingSetList)
	http.HandleFunc("/", h.GETReadingType)
	http.HandleFunc("/", h.GETUsagePoint)
	http.HandleFunc("/", h.GETUsagePointList)
	http.HandleFunc("/", h.HEADMeterReading)
	http.HandleFunc("/", h.HEADMeterReadingList)
	http.HandleFunc("/", h.HEADReading)
	http.HandleFunc("/", h.HEADReadingList)
	http.HandleFunc("/", h.HEADReadingSet)
	http.HandleFunc("/", h.HEADReadingSetList)
	http.HandleFunc("/", h.HEADReadingType)
	http.HandleFunc("/", h.HEADUsagePoint)
	http.HandleFunc("/", h.HEADUsagePointList)
	http.HandleFunc("/", h.POSTMeterReading)
	http.HandleFunc("/", h.POSTMeterReadingList)
	http.HandleFunc("/", h.POSTReading)
	http.HandleFunc("/", h.POSTReadingList)
	http.HandleFunc("/", h.POSTReadingSet)
	http.HandleFunc("/", h.POSTReadingSetList)
	http.HandleFunc("/", h.POSTReadingType)
	http.HandleFunc("/", h.POSTUsagePoint)
	http.HandleFunc("/", h.POSTUsagePointList)
	http.HandleFunc("/", h.PUTMeterReading)
	http.HandleFunc("/", h.PUTMeterReadingList)
	http.HandleFunc("/", h.PUTReading)
	http.HandleFunc("/", h.PUTReadingList)
	http.HandleFunc("/", h.PUTReadingSet)
	http.HandleFunc("/", h.PUTReadingSetList)
	http.HandleFunc("/", h.PUTReadingType)
	http.HandleFunc("/", h.PUTUsagePoint)
	http.HandleFunc("/", h.PUTUsagePointList)

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
