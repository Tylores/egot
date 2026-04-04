package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/MUP/handler"
"github.com/Tylores/egot/internal/MUP/repository/memory"
"github.com/Tylores/egot/internal/routes"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      routes.MUP,
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)
	http.Handle("GET /mup", http.HandlerFunc(h.GETMirrorUsagePointList))
	http.Handle("HEAD /mup", http.HandlerFunc(h.HEADMirrorUsagePointList))
	http.Handle("PUT /mup", http.HandlerFunc(h.PUTMirrorUsagePointList))
	http.Handle("POST /mup", http.HandlerFunc(h.POSTMirrorUsagePointList))
	http.Handle("DELETE /mup", http.HandlerFunc(h.DELETEMirrorUsagePointList))
	http.Handle("GET /mup/{id1}", http.HandlerFunc(h.GETMirrorUsagePoint))
	http.Handle("HEAD /mup/{id1}", http.HandlerFunc(h.HEADMirrorUsagePoint))
	http.Handle("PUT /mup/{id1}", http.HandlerFunc(h.PUTMirrorUsagePoint))
	http.Handle("POST /mup/{id1}", http.HandlerFunc(h.POSTMirrorUsagePoint))
	http.Handle("DELETE /mup/{id1}", http.HandlerFunc(h.DELETEMirrorUsagePoint))

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
