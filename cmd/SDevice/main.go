package main

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/SDevice/handler"
	"github.com/Tylores/egot/internal/SDevice/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.SDevice,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository()

	h := handler.NewHandler(repo)
	http.Handle("GET /sdev", http.HandlerFunc(h.GETSelfDevice))
	http.Handle("HEAD /sdev", http.HandlerFunc(h.HEADSelfDevice))
	http.Handle("PUT /sdev", http.HandlerFunc(h.PUTSelfDevice))
	http.Handle("POST /sdev", http.HandlerFunc(h.POSTSelfDevice))
	http.Handle("DELETE /sdev", http.HandlerFunc(h.DELETESelfDevice))

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
