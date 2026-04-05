package main

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/UPT/handler"
	"github.com/Tylores/egot/internal/UPT/repository/memory"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      "egot.internal.com:8024",
		TLSConfig: cfg,
	}

	repo := memory.NewRepository()

	h := handler.NewHandler(repo)

	// Register handler methods
	http.Handle("DELETE /upt", http.HandlerFunc(h.DELETEUsagePointList))
	http.Handle("GET /upt", http.HandlerFunc(h.GETUsagePointList))
	http.Handle("HEAD /upt", http.HandlerFunc(h.HEADUsagePointList))
	http.Handle("POST /upt", http.HandlerFunc(h.POSTUsagePointList))
	http.Handle("PUT /upt", http.HandlerFunc(h.PUTUsagePointList))
	http.Handle("DELETE /upt/{id1}", http.HandlerFunc(h.DELETEUsagePoint))
	http.Handle("GET /upt/{id1}", http.HandlerFunc(h.GETUsagePoint))
	http.Handle("HEAD /upt/{id1}", http.HandlerFunc(h.HEADUsagePoint))
	http.Handle("POST /upt/{id1}", http.HandlerFunc(h.POSTUsagePoint))
	http.Handle("PUT /upt/{id1}", http.HandlerFunc(h.PUTUsagePoint))
	http.Handle("DELETE /upt/{id1}/mr", http.HandlerFunc(h.DELETEMeterReadingList))
	http.Handle("GET /upt/{id1}/mr", http.HandlerFunc(h.GETMeterReadingList))
	http.Handle("HEAD /upt/{id1}/mr", http.HandlerFunc(h.HEADMeterReadingList))
	http.Handle("POST /upt/{id1}/mr", http.HandlerFunc(h.POSTMeterReadingList))
	http.Handle("PUT /upt/{id1}/mr", http.HandlerFunc(h.PUTMeterReadingList))
	http.Handle("DELETE /upt/{id1}/mr/{id2}", http.HandlerFunc(h.DELETEMeterReading))
	http.Handle("GET /upt/{id1}/mr/{id2}", http.HandlerFunc(h.GETMeterReading))
	http.Handle("HEAD /upt/{id1}/mr/{id2}", http.HandlerFunc(h.HEADMeterReading))
	http.Handle("POST /upt/{id1}/mr/{id2}", http.HandlerFunc(h.POSTMeterReading))
	http.Handle("PUT /upt/{id1}/mr/{id2}", http.HandlerFunc(h.PUTMeterReading))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.DELETEReadingSetList))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.GETReadingSetList))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.HEADReadingSetList))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.POSTReadingSetList))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs", http.HandlerFunc(h.PUTReadingSetList))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.DELETEReadingSet))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.GETReadingSet))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.HEADReadingSet))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.POSTReadingSet))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}", http.HandlerFunc(h.PUTReadingSet))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.DELETEReadingList))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.GETReadingList))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.HEADReadingList))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.POSTReadingList))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}/r", http.HandlerFunc(h.PUTReadingList))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.DELETEReading))
	http.Handle("GET /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.GETReading))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.HEADReading))
	http.Handle("POST /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.POSTReading))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rs/{id3}/r/{id4}", http.HandlerFunc(h.PUTReading))
	http.Handle("DELETE /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.DELETEReadingType))
	http.Handle("GET /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.GETReadingType))
	http.Handle("HEAD /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.HEADReadingType))
	http.Handle("POST /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.POSTReadingType))
	http.Handle("PUT /upt/{id1}/mr/{id2}/rt", http.HandlerFunc(h.PUTReadingType))
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
