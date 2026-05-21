package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/Notify/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.Notify,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs

	repo := store.New(filepath.Join("data", "Notify.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /ntfy", http.HandlerFunc(h.GETNotificationList))
	http.Handle("HEAD /ntfy", http.HandlerFunc(h.HEADNotificationList))
	http.Handle("PUT /ntfy", http.HandlerFunc(h.PUTNotificationList))
	http.Handle("POST /ntfy", http.HandlerFunc(h.POSTNotificationList))
	http.Handle("DELETE /ntfy", http.HandlerFunc(h.DELETENotificationList))
	http.Handle("GET /ntfy/{id1}", http.HandlerFunc(h.GETNotification))
	http.Handle("HEAD /ntfy/{id1}", http.HandlerFunc(h.HEADNotification))
	http.Handle("PUT /ntfy/{id1}", http.HandlerFunc(h.PUTNotification))
	http.Handle("POST /ntfy/{id1}", http.HandlerFunc(h.POSTNotification))
	http.Handle("DELETE /ntfy/{id1}", http.HandlerFunc(h.DELETENotification))

	log.Printf("Starting Notify on %s", routes.Notify)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
