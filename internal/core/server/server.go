package server

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/core/handler"
	"github.com/Tylores/egot/internal/core/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
	"github.com/Tylores/egot/sep/uri"
)

func AddRoutes(h *handler.Handler) {
	http.Handle("GET "+uri.DeviceCapability, http.HandlerFunc(h.GetDeviceCapability))
	http.Handle("GET "+uri.Time, http.HandlerFunc(h.GetTime))
	http.Handle("GET "+uri.EndDeviceList, http.HandlerFunc(h.GetEndDevices))
	http.Handle("GET "+uri.EndDevice, http.HandlerFunc(h.GetEndDevice))
	http.Handle("GET "+uri.Registration, http.HandlerFunc(h.GetRegistration))
}

func ServeHTTP(entities memory.Entity) {
	server := http.Server{
		Addr: routes.Core,
	}

	repo := memory.NewRepository(entities)
	h := handler.NewHandler(repo)
	http.Handle("GET "+uri.DeviceCapability, http.HandlerFunc(h.GetDeviceCapability))
	http.Handle("GET "+uri.Time, http.HandlerFunc(h.GetTime))
	http.Handle("GET "+uri.EndDeviceList, http.HandlerFunc(h.GetEndDevices))
	http.Handle("GET "+uri.EndDevice, http.HandlerFunc(h.GetEndDevice))
	http.Handle("GET "+uri.Registration, http.HandlerFunc(h.GetRegistration))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func ServeHTTPS(entities memory.Entity) {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.Core,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository(entities)
	repo.InitRepository("./ssl")

	h := handler.NewHandler(repo)
	AddRoutes(h)

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
