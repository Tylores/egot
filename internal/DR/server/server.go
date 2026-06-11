package server

import (
	"time"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/DR/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func AddRoutes(h *handler.Handler) {
	// Add routes from WADL specification
}

func ServeHTTPS() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:              routes.DR,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig:         cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	defer reg.Close()

	repo := store.New(filepath.Join("data", "DR.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
	AddRoutes(h)

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
