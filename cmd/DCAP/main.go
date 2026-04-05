package main

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/DCAP/handler"
	"github.com/Tylores/egot/internal/DCAP/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.DCAP,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository()

	h := handler.NewHandler(repo)
	http.Handle("GET /dcap", http.HandlerFunc(h.GETDeviceCapability))
	http.Handle("HEAD /dcap", http.HandlerFunc(h.HEADDeviceCapability))
	http.Handle("PUT /dcap", http.HandlerFunc(h.PUTDeviceCapability))
	http.Handle("POST /dcap", http.HandlerFunc(h.POSTDeviceCapability))
	http.Handle("DELETE /dcap", http.HandlerFunc(h.DELETEDeviceCapability))

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
