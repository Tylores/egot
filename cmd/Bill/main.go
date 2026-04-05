package main

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/Bill/handler"
	"github.com/Tylores/egot/internal/Bill/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.Bill,
		TLSConfig: cfg,
	}

	repo := memory.NewRepository()

	h := handler.NewHandler(repo)
	http.Handle("GET /bill", http.HandlerFunc(h.GETCustomerAccountList))
	http.Handle("HEAD /bill", http.HandlerFunc(h.HEADCustomerAccountList))
	http.Handle("PUT /bill", http.HandlerFunc(h.PUTCustomerAccountList))
	http.Handle("POST /bill", http.HandlerFunc(h.POSTCustomerAccountList))
	http.Handle("DELETE /bill", http.HandlerFunc(h.DELETECustomerAccountList))
	http.Handle("GET /bill/{id1}", http.HandlerFunc(h.GETCustomerAccount))
	http.Handle("HEAD /bill/{id1}", http.HandlerFunc(h.HEADCustomerAccount))
	http.Handle("PUT /bill/{id1}", http.HandlerFunc(h.PUTCustomerAccount))
	http.Handle("POST /bill/{id1}", http.HandlerFunc(h.POSTCustomerAccount))
	http.Handle("DELETE /bill/{id1}", http.HandlerFunc(h.DELETECustomerAccount))
	http.Handle("GET /bill/{id1}/ca", http.HandlerFunc(h.GETCustomerAgreementList))
	http.Handle("HEAD /bill/{id1}/ca", http.HandlerFunc(h.HEADCustomerAgreementList))
	http.Handle("PUT /bill/{id1}/ca", http.HandlerFunc(h.PUTCustomerAgreementList))
	http.Handle("POST /bill/{id1}/ca", http.HandlerFunc(h.POSTCustomerAgreementList))
	http.Handle("DELETE /bill/{id1}/ca", http.HandlerFunc(h.DELETECustomerAgreementList))
	http.Handle("GET /bill/{id1}/ca/{id2}", http.HandlerFunc(h.GETCustomerAgreement))
	http.Handle("HEAD /bill/{id1}/ca/{id2}", http.HandlerFunc(h.HEADCustomerAgreement))
	http.Handle("PUT /bill/{id1}/ca/{id2}", http.HandlerFunc(h.PUTCustomerAgreement))
	http.Handle("POST /bill/{id1}/ca/{id2}", http.HandlerFunc(h.POSTCustomerAgreement))
	http.Handle("DELETE /bill/{id1}/ca/{id2}", http.HandlerFunc(h.DELETECustomerAgreement))
	http.Handle("GET /bill/{id1}/ca/{id2}/actbp", http.HandlerFunc(h.GETActiveBillingPeriodList))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/actbp", http.HandlerFunc(h.HEADActiveBillingPeriodList))
	http.Handle("PUT /bill/{id1}/ca/{id2}/actbp", http.HandlerFunc(h.PUTActiveBillingPeriodList))
	http.Handle("POST /bill/{id1}/ca/{id2}/actbp", http.HandlerFunc(h.POSTActiveBillingPeriodList))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/actbp", http.HandlerFunc(h.DELETEActiveBillingPeriodList))
	http.Handle("GET /bill/{id1}/ca/{id2}/bp", http.HandlerFunc(h.GETBillingPeriodList))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/bp", http.HandlerFunc(h.HEADBillingPeriodList))
	http.Handle("PUT /bill/{id1}/ca/{id2}/bp", http.HandlerFunc(h.PUTBillingPeriodList))
	http.Handle("POST /bill/{id1}/ca/{id2}/bp", http.HandlerFunc(h.POSTBillingPeriodList))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/bp", http.HandlerFunc(h.DELETEBillingPeriodList))
	http.Handle("GET /bill/{id1}/ca/{id2}/bp/{id3}", http.HandlerFunc(h.GETBillingPeriod))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/bp/{id3}", http.HandlerFunc(h.HEADBillingPeriod))
	http.Handle("PUT /bill/{id1}/ca/{id2}/bp/{id3}", http.HandlerFunc(h.PUTBillingPeriod))
	http.Handle("POST /bill/{id1}/ca/{id2}/bp/{id3}", http.HandlerFunc(h.POSTBillingPeriod))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/bp/{id3}", http.HandlerFunc(h.DELETEBillingPeriod))
	http.Handle("GET /bill/{id1}/ca/{id2}/pro", http.HandlerFunc(h.GETProjectionReadingList))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/pro", http.HandlerFunc(h.HEADProjectionReadingList))
	http.Handle("PUT /bill/{id1}/ca/{id2}/pro", http.HandlerFunc(h.PUTProjectionReadingList))
	http.Handle("POST /bill/{id1}/ca/{id2}/pro", http.HandlerFunc(h.POSTProjectionReadingList))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/pro", http.HandlerFunc(h.DELETEProjectionReadingList))
	http.Handle("GET /bill/{id1}/ca/{id2}/pro/{id3}", http.HandlerFunc(h.GETProjectionReading))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/pro/{id3}", http.HandlerFunc(h.HEADProjectionReading))
	http.Handle("PUT /bill/{id1}/ca/{id2}/pro/{id3}", http.HandlerFunc(h.PUTProjectionReading))
	http.Handle("POST /bill/{id1}/ca/{id2}/pro/{id3}", http.HandlerFunc(h.POSTProjectionReading))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/pro/{id3}", http.HandlerFunc(h.DELETEProjectionReading))
	http.Handle("GET /bill/{id1}/ca/{id2}/tar", http.HandlerFunc(h.GETTargetReadingList))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/tar", http.HandlerFunc(h.HEADTargetReadingList))
	http.Handle("PUT /bill/{id1}/ca/{id2}/tar", http.HandlerFunc(h.PUTTargetReadingList))
	http.Handle("POST /bill/{id1}/ca/{id2}/tar", http.HandlerFunc(h.POSTTargetReadingList))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/tar", http.HandlerFunc(h.DELETETargetReadingList))
	http.Handle("GET /bill/{id1}/ca/{id2}/tar/{id3}", http.HandlerFunc(h.GETTargetReading))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/tar/{id3}", http.HandlerFunc(h.HEADTargetReading))
	http.Handle("PUT /bill/{id1}/ca/{id2}/tar/{id3}", http.HandlerFunc(h.PUTTargetReading))
	http.Handle("POST /bill/{id1}/ca/{id2}/tar/{id3}", http.HandlerFunc(h.POSTTargetReading))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/tar/{id3}", http.HandlerFunc(h.DELETETargetReading))
	http.Handle("GET /bill/{id1}/ca/{id2}/ver", http.HandlerFunc(h.GETHistoricalReadingList))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/ver", http.HandlerFunc(h.HEADHistoricalReadingList))
	http.Handle("PUT /bill/{id1}/ca/{id2}/ver", http.HandlerFunc(h.PUTHistoricalReadingList))
	http.Handle("POST /bill/{id1}/ca/{id2}/ver", http.HandlerFunc(h.POSTHistoricalReadingList))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/ver", http.HandlerFunc(h.DELETEHistoricalReadingList))
	http.Handle("GET /bill/{id1}/ca/{id2}/ver/{id3}", http.HandlerFunc(h.GETHistoricalReading))
	http.Handle("HEAD /bill/{id1}/ca/{id2}/ver/{id3}", http.HandlerFunc(h.HEADHistoricalReading))
	http.Handle("PUT /bill/{id1}/ca/{id2}/ver/{id3}", http.HandlerFunc(h.PUTHistoricalReading))
	http.Handle("POST /bill/{id1}/ca/{id2}/ver/{id3}", http.HandlerFunc(h.POSTHistoricalReading))
	http.Handle("DELETE /bill/{id1}/ca/{id2}/ver/{id3}", http.HandlerFunc(h.DELETEHistoricalReading))
	http.Handle("GET /bill/{id1}/ss", http.HandlerFunc(h.GETServiceSupplier))
	http.Handle("HEAD /bill/{id1}/ss", http.HandlerFunc(h.HEADServiceSupplier))
	http.Handle("PUT /bill/{id1}/ss", http.HandlerFunc(h.PUTServiceSupplier))
	http.Handle("POST /bill/{id1}/ss", http.HandlerFunc(h.POSTServiceSupplier))
	http.Handle("DELETE /bill/{id1}/ss", http.HandlerFunc(h.DELETEServiceSupplier))

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
