package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/sep/uri"
)

func main() {
	caCert, _ := os.ReadFile("./ssl/ca.crt")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, _ := tls.LoadX509KeyPair("./ssl/client.crt", "./ssl/client.key")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get("https://" + routes.core + uri.DeviceCapability)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", dump)

	resp, err = client.Head("https://" + routes.core + uri.DeviceCapability)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", dump)

	resp, err = client.Get("https://" + routes.core + uri.Time)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", dump)

	resp, err = client.Head("https://" + routes.core + uri.Time)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", dump)
}
