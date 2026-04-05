package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/TimeOfUse/handler"
"github.com/Tylores/egot/internal/TimeOfUse/repository/memory"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      "egot.internal.com:8023",
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)

// Register handler methods
	http.Handle("DELETE /tm", http.HandlerFunc(h.DELETETime))
	http.Handle("GET /tm", http.HandlerFunc(h.GETTime))
	http.Handle("HEAD /tm", http.HandlerFunc(h.HEADTime))
	http.Handle("POST /tm", http.HandlerFunc(h.POSTTime))
	http.Handle("PUT /tm", http.HandlerFunc(h.PUTTime))
err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
