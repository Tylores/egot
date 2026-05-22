package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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
		Handler:   tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	defer reg.Close()
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "Notify.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting Notify on %s", routes.Notify)
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down Notify server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
}
