package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
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
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  pool,
	}, nil
}
