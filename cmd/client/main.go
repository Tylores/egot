package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/sep/uri"
)

func main() {
	caCert, err := os.ReadFile("./ssl/ca.crt")
	if err != nil {
		panic(fmt.Errorf("failed to read CA certificate: %w", err))
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		panic("failed to parse CA certificate")
	}

	cert, err := tls.LoadX509KeyPair("./ssl/client.crt", "./ssl/client.key")
	if err != nil {
		panic(fmt.Errorf("failed to load client key pair: %w", err))
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
		Timeout: 10 * time.Second,
	}

	call := func(method, path string) {
		req, err := http.NewRequest(method, "https://"+routes.DCAP+path, nil)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", dump)
	}

	call("GET", uri.DeviceCapability)
	call("HEAD", uri.DeviceCapability)
	call("GET", uri.Time)
	call("HEAD", uri.Time)
}
