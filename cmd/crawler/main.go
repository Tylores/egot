//go:build xmlvalidate

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/sep"
	"github.com/Tylores/egot/sep/uri"
	"github.com/terminalstatic/go-xsd-validate"
)

var (
	visited   = make(map[string]bool)
	visitedMu sync.Mutex
)

func shouldVisit(href string) bool {
	visitedMu.Lock()
	defer visitedMu.Unlock()
	if visited[href] {
		return false
	}
	visited[href] = true
	return true
}

func validate(body []byte) {
	xsdvalidate.Init()
	defer xsdvalidate.Cleanup()

	xsd_handler, err := xsdvalidate.NewXsdHandlerUrl("internal/sep/sep.xsd", xsdvalidate.ParsErrDefault)
	if err != nil {
		log.Fatal(err)
	}
	defer xsd_handler.Free()

	xml_handler, err := xsdvalidate.NewXmlHandlerMem(body, xsdvalidate.ValidErrDefault)
	if err != nil {
		log.Println(err)
		return
	}
	defer xml_handler.Free()

	err = xsd_handler.ValidateMem(body, xsdvalidate.ValidErrDefault)
	if err != nil {
		log.Println(string(body))
		switch err.(type) {
		case xsdvalidate.ValidationError:
			log.Println(err)
			log.Printf("Error in line: %d\n", err.(xsdvalidate.ValidationError).Errors[0].Line)
			log.Println(err.(xsdvalidate.ValidationError).Errors[0].Message)
		default:
			log.Println(err)
		}
	}
}

func checkHead(client *http.Client, href string) int {
	url := "https://" + routes.Gateway + href
	resp, err := client.Head(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	return resp.StatusCode
}

func checkGet(client *http.Client, href string) int {
	url := "https://" + routes.Gateway + href
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		log.Printf("Validating %s", href)
		validate(body)
	}

	re := regexp.MustCompile(`href="(.*?)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		new_href := string(match[1])
		if new_href != href {
			checkAccess(client, match[1])
		}
	}

	return resp.StatusCode
}

func checkPut(client *http.Client, href string) int {
	url := "https://" + routes.Gateway + href
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", sep.ContentType)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	return resp.StatusCode
}

func checkPost(client *http.Client, href string) int {
	url := "https://" + routes.Gateway + href
	resp, err := client.Post(url, sep.ContentType, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	return resp.StatusCode
}

func checkDelete(client *http.Client, href string) int {
	url := "https://" + routes.Gateway + href
	req, err := http.NewRequest("DELETE", url, nil)
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
	return resp.StatusCode
}

func checkAccess(client *http.Client, href string) {
	if !shouldVisit(href) {
		return
	}
	head_code := checkHead(client, href)
	get_code := checkGet(client, href)
	put_code := checkPut(client, href)
	post_code := checkPost(client, href)
	delete_code := checkDelete(client, href)
	log.Printf("%s HEAD:%d GET:%d POST:%d PUT:%d DELETE:%d", href, head_code, get_code, put_code, post_code, delete_code)
}

func main() {
	caCert, err := os.ReadFile("./ssl/ca.crt")
	if err != nil {
		log.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to parse CA certificate")
	}

	cert, err := tls.LoadX509KeyPair("./ssl/client.crt", "./ssl/client.key")
	if err != nil {
		log.Fatalf("failed to load client key pair: %v", err)
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

	checkAccess(client, uri.DeviceCapability)
}
