package servers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func MTLSConfig(caCert, cert, key string) (*tls.Config, error) {
	loadedServerCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	loadedCACert, err := os.ReadFile(caCert)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(loadedCACert) {
		return nil, fmt.Errorf("cannot add CA cert in the pool")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{loadedServerCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}
	return tlsConfig, nil
}
