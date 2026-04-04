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
	http.HandleFunc("/", h.DELETETime)
	http.HandleFunc("/", h.GETTime)
	http.HandleFunc("/", h.HEADTime)
	http.HandleFunc("/", h.POSTTime)
	http.HandleFunc("/", h.PUTTime)

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
