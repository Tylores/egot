package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	server := http.Server{
		Addr:      ":4443",
		TLSConfig: cfg,
	}
	err := server.ListenAndServeTLS("./ssl/srv.crt", "./ssl/srv.key")
	if err != nil {
		log.Fatal(err)
	}
}
