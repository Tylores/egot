package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// NewServerConfig returns a TLS configuration for mTLS servers.
// sslDir must contain ca.crt, which is used to verify client certificates.
func NewServerConfig(sslDir string) (*tls.Config, error) {
	ca, err := os.ReadFile(sslDir + "/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("read CA cert from %s: %w", sslDir, err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("failed to parse CA certificate from %s/ca.crt", sslDir)
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  pool,
	}, nil
}

// CertHeaderMiddleware extracts the client certificate from the X-SSL-Client-Cert header
// and injects it into req.TLS.PeerCertificates, making it transparent to the handlers.
func CertHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.TLS == nil {
			req.TLS = &tls.ConnectionState{}
		}

		if len(req.TLS.PeerCertificates) == 0 {
			if headerVal := req.Header.Get("X-SSL-Client-Cert"); headerVal != "" {
				unescaped, err := url.PathUnescape(headerVal)
				if err == nil && unescaped != "" {
					block, _ := pem.Decode([]byte(unescaped))
					if block != nil && block.Type == "CERTIFICATE" {
						cert, err := x509.ParseCertificate(block.Bytes)
						if err == nil {
							req.TLS.PeerCertificates = []*x509.Certificate{cert}
						}
					}
				}
			}
		}

		next.ServeHTTP(w, req)
	})
}
