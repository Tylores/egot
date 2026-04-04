package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/File/handler"
"github.com/Tylores/egot/internal/File/repository/memory"
"github.com/Tylores/egot/internal/routes"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      routes.File,
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)
	http.Handle("GET /file", http.HandlerFunc(h.GETFileList))
	http.Handle("HEAD /file", http.HandlerFunc(h.HEADFileList))
	http.Handle("PUT /file", http.HandlerFunc(h.PUTFileList))
	http.Handle("POST /file", http.HandlerFunc(h.POSTFileList))
	http.Handle("DELETE /file", http.HandlerFunc(h.DELETEFileList))
	http.Handle("GET /file/{id1}", http.HandlerFunc(h.GETFile))
	http.Handle("HEAD /file/{id1}", http.HandlerFunc(h.HEADFile))
	http.Handle("PUT /file/{id1}", http.HandlerFunc(h.PUTFile))
	http.Handle("POST /file/{id1}", http.HandlerFunc(h.POSTFile))
	http.Handle("DELETE /file/{id1}", http.HandlerFunc(h.DELETEFile))

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
