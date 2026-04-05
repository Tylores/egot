package main

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/DERP/handler"
	"github.com/Tylores/egot/internal/DERP/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.DERP,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository()

	h := handler.NewHandler(repo)
	http.Handle("GET /derp", http.HandlerFunc(h.GETDERProgramList))
	http.Handle("HEAD /derp", http.HandlerFunc(h.HEADDERProgramList))
	http.Handle("PUT /derp", http.HandlerFunc(h.PUTDERProgramList))
	http.Handle("POST /derp", http.HandlerFunc(h.POSTDERProgramList))
	http.Handle("DELETE /derp", http.HandlerFunc(h.DELETEDERProgramList))
	http.Handle("GET /derp/{id1}", http.HandlerFunc(h.GETDERProgram))
	http.Handle("HEAD /derp/{id1}", http.HandlerFunc(h.HEADDERProgram))
	http.Handle("PUT /derp/{id1}", http.HandlerFunc(h.PUTDERProgram))
	http.Handle("POST /derp/{id1}", http.HandlerFunc(h.POSTDERProgram))
	http.Handle("DELETE /derp/{id1}", http.HandlerFunc(h.DELETEDERProgram))
	http.Handle("GET /derp/{id1}/actderc", http.HandlerFunc(h.GETActiveDERControlList))
	http.Handle("HEAD /derp/{id1}/actderc", http.HandlerFunc(h.HEADActiveDERControlList))
	http.Handle("PUT /derp/{id1}/actderc", http.HandlerFunc(h.PUTActiveDERControlList))
	http.Handle("POST /derp/{id1}/actderc", http.HandlerFunc(h.POSTActiveDERControlList))
	http.Handle("DELETE /derp/{id1}/actderc", http.HandlerFunc(h.DELETEActiveDERControlList))
	http.Handle("GET /derp/{id1}/derc", http.HandlerFunc(h.GETDERControlList))
	http.Handle("HEAD /derp/{id1}/derc", http.HandlerFunc(h.HEADDERControlList))
	http.Handle("PUT /derp/{id1}/derc", http.HandlerFunc(h.PUTDERControlList))
	http.Handle("POST /derp/{id1}/derc", http.HandlerFunc(h.POSTDERControlList))
	http.Handle("DELETE /derp/{id1}/derc", http.HandlerFunc(h.DELETEDERControlList))
	http.Handle("GET /derp/{id1}/derc/{id2}", http.HandlerFunc(h.GETDERControl))
	http.Handle("HEAD /derp/{id1}/derc/{id2}", http.HandlerFunc(h.HEADDERControl))
	http.Handle("PUT /derp/{id1}/derc/{id2}", http.HandlerFunc(h.PUTDERControl))
	http.Handle("POST /derp/{id1}/derc/{id2}", http.HandlerFunc(h.POSTDERControl))
	http.Handle("DELETE /derp/{id1}/derc/{id2}", http.HandlerFunc(h.DELETEDERControl))
	http.Handle("GET /derp/{id1}/dderc", http.HandlerFunc(h.GETDefaultDERControl))
	http.Handle("HEAD /derp/{id1}/dderc", http.HandlerFunc(h.HEADDefaultDERControl))
	http.Handle("PUT /derp/{id1}/dderc", http.HandlerFunc(h.PUTDefaultDERControl))
	http.Handle("POST /derp/{id1}/dderc", http.HandlerFunc(h.POSTDefaultDERControl))
	http.Handle("DELETE /derp/{id1}/dderc", http.HandlerFunc(h.DELETEDefaultDERControl))
	http.Handle("GET /derp/{id1}/dc", http.HandlerFunc(h.GETDERCurveList))
	http.Handle("HEAD /derp/{id1}/dc", http.HandlerFunc(h.HEADDERCurveList))
	http.Handle("PUT /derp/{id1}/dc", http.HandlerFunc(h.PUTDERCurveList))
	http.Handle("POST /derp/{id1}/dc", http.HandlerFunc(h.POSTDERCurveList))
	http.Handle("DELETE /derp/{id1}/dc", http.HandlerFunc(h.DELETEDERCurveList))
	http.Handle("GET /derp/{id1}/dc/{id2}", http.HandlerFunc(h.GETDERCurve))
	http.Handle("HEAD /derp/{id1}/dc/{id2}", http.HandlerFunc(h.HEADDERCurve))
	http.Handle("PUT /derp/{id1}/dc/{id2}", http.HandlerFunc(h.PUTDERCurve))
	http.Handle("POST /derp/{id1}/dc/{id2}", http.HandlerFunc(h.POSTDERCurve))
	http.Handle("DELETE /derp/{id1}/dc/{id2}", http.HandlerFunc(h.DELETEDERCurve))

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
