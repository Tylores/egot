package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/core/handler"
	"github.com/Tylores/egot/internal/core/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/sep/uri"
)

const MAX_ENTITIES memory.Entity = 10

func main() {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	server := http.Server{
		Addr:      routes.Core,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository(MAX_ENTITIES)
	repo.InitRepository("./ssl")

	h := handler.NewHandler(repo)
	http.Handle("GET "+uri.DeviceCapability, http.HandlerFunc(h.GetDeviceCapability))
	http.Handle("GET "+uri.Time, http.HandlerFunc(h.GetTime))
	http.Handle("GET "+uri.EndDeviceList, http.HandlerFunc(h.GetEndDevices))
	http.Handle("GET "+uri.EndDevice, http.HandlerFunc(h.GetEndDevice))
	http.Handle("GET "+uri.Registration, http.HandlerFunc(h.GetRegistration))

	err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
