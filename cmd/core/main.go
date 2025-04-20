package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/core/handler"
	"github.com/Tylores/egot/internal/core/repository/memory"
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
	repo := memory.NewRepository(1)
	h := handler.NewHandler(repo)
	http.Handle("/dcap", http.HandlerFunc(h.GetDeviceCapability))
	err := server.ListenAndServeTLS("./ssl/srv.crt", "./ssl/srv.key")
	if err != nil {
		log.Fatal(err)
	}
}
