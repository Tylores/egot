package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/Messaging/handler"
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
		Addr:      routes.Messaging,
		TLSConfig: cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "Messaging.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(repo, reg)
	http.Handle("GET /msg", http.HandlerFunc(h.GETMessagingProgramList))
	http.Handle("HEAD /msg", http.HandlerFunc(h.HEADMessagingProgramList))
	http.Handle("PUT /msg", http.HandlerFunc(h.PUTMessagingProgramList))
	http.Handle("POST /msg", http.HandlerFunc(h.POSTMessagingProgramList))
	http.Handle("DELETE /msg", http.HandlerFunc(h.DELETEMessagingProgramList))
	http.Handle("GET /msg/{id1}", http.HandlerFunc(h.GETMessagingProgram))
	http.Handle("HEAD /msg/{id1}", http.HandlerFunc(h.HEADMessagingProgram))
	http.Handle("PUT /msg/{id1}", http.HandlerFunc(h.PUTMessagingProgram))
	http.Handle("POST /msg/{id1}", http.HandlerFunc(h.POSTMessagingProgram))
	http.Handle("DELETE /msg/{id1}", http.HandlerFunc(h.DELETEMessagingProgram))
	http.Handle("GET /msg/{id1}/acttxt", http.HandlerFunc(h.GETActiveTextMessageList))
	http.Handle("HEAD /msg/{id1}/acttxt", http.HandlerFunc(h.HEADActiveTextMessageList))
	http.Handle("PUT /msg/{id1}/acttxt", http.HandlerFunc(h.PUTActiveTextMessageList))
	http.Handle("POST /msg/{id1}/acttxt", http.HandlerFunc(h.POSTActiveTextMessageList))
	http.Handle("DELETE /msg/{id1}/acttxt", http.HandlerFunc(h.DELETEActiveTextMessageList))
	http.Handle("GET /msg/{id1}/txt", http.HandlerFunc(h.GETTextMessageList))
	http.Handle("HEAD /msg/{id1}/txt", http.HandlerFunc(h.HEADTextMessageList))
	http.Handle("PUT /msg/{id1}/txt", http.HandlerFunc(h.PUTTextMessageList))
	http.Handle("POST /msg/{id1}/txt", http.HandlerFunc(h.POSTTextMessageList))
	http.Handle("DELETE /msg/{id1}/txt", http.HandlerFunc(h.DELETETextMessageList))
	http.Handle("GET /msg/{id1}/txt/{id2}", http.HandlerFunc(h.GETTextMessage))
	http.Handle("HEAD /msg/{id1}/txt/{id2}", http.HandlerFunc(h.HEADTextMessage))
	http.Handle("PUT /msg/{id1}/txt/{id2}", http.HandlerFunc(h.PUTTextMessage))
	http.Handle("POST /msg/{id1}/txt/{id2}", http.HandlerFunc(h.POSTTextMessage))
	http.Handle("DELETE /msg/{id1}/txt/{id2}", http.HandlerFunc(h.DELETETextMessage))

	log.Printf("Starting Messaging on %s", routes.Messaging)
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
