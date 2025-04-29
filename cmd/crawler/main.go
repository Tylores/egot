package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/sep"
	"github.com/Tylores/egot/internal/uri"
)

func CheckHead(client *http.Client, href string) int {
	url := "https://" + routes.Core + href
	resp, err := client.Head(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func CheckGet(client *http.Client, href string) int {
	url := "https://" + routes.Core + href
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))

	re := regexp.MustCompile(`href="(.*?)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		new_href := string(match[1])
		if new_href != href {
			CheckAccess(client, match[1])
		}
	}

	return resp.StatusCode
}

func CheckPut(client *http.Client, href string) int {
	url := "https://" + routes.Core + href
	req, err := http.NewRequest("PUT", url, nil)
	req.Header.Set("Content-Type", sep.ContentType)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func CheckPost(client *http.Client, href string) int {
	url := "https://" + routes.Core + href
	resp, err := client.Post(url, sep.ContentType, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func CheckDelete(client *http.Client, href string) int {
	url := "https://" + routes.Core + href
	req, err := http.NewRequest("DELETE", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func CheckAccess(client *http.Client, href string) {
	head_code := CheckHead(client, href)
	get_code := CheckGet(client, href)
	put_code := CheckPut(client, href)
	post_code := CheckPost(client, href)
	delete_code := CheckDelete(client, href)
	log.Printf("%s HEAD:%d GET:%d POST:%d PUT:%d DELETE:%d", href, head_code, get_code, put_code, post_code, delete_code)
}

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

	CheckAccess(client, uri.DeviceCapability)
}
