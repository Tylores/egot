package server

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/rsps/handler"
	"github.com/Tylores/egot/internal/rsps/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func AddRoutes(h *handler.Handler) {
	// Add routes from WADL specification
	// Example: http.Handle("GET /resource", http.HandlerFunc(h.GetResource))
}

func ServeHTTPS(entities memory.Entity) {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.Rsps,
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
