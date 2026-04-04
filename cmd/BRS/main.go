package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/BRS/handler"
"github.com/Tylores/egot/internal/BRS/repository/memory"
"github.com/Tylores/egot/internal/routes"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      routes.BRS,
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)
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

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
