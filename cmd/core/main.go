package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/core/handler"
	"github.com/Tylores/egot/internal/core/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/uri"
)

func main() {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	server := http.Server{
		Addr:      routes.Core,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository(10)
	repo.InitRepository("./ssl")

	h := handler.NewHandler(repo)
	http.Handle(uri.DeviceCapability, http.HandlerFunc(h.GetDeviceCapability))
	http.Handle(uri.Time, http.HandlerFunc(h.GetTime))

	err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
