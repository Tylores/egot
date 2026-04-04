package main

import (
"crypto/tls"
"log"
"net/http"

"github.com/Tylores/egot/internal/Notify/handler"
"github.com/Tylores/egot/internal/Notify/repository/memory"
"github.com/Tylores/egot/internal/routes"
)

func main() {
cfg := &tls.Config{
MinVersion: tls.VersionTLS12,
ClientAuth: tls.RequireAndVerifyClientCert,
}
server := http.Server{
Addr:      routes.Notify,
TLSConfig: cfg,
}

repo := memory.NewRepository()

h := handler.NewHandler(repo)
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

err := server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
if err != nil {
log.Fatal(err)
}
}
