package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/ssl"
)

func NewTLSConfig(cfg config.Config) (*tls.Config, error) {

	cert, err := ssl.LoadCertificate(cfg.SSL.Server.Cert)
	if err != nil {
		return nil, fmt.Errorf("cannot load server cert: %w", err)
	}

	key, err := ssl.LoadPrivateKey(cfg.SSL.Server.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot load server key: %w", err)
	}

	caCert, err := ssl.LoadCertificate(cfg.SSL.CA.Cert)
	if err != nil {
		return nil, fmt.Errorf("cannot load CA cert: %w", err)
	}

	tlsCert := tls.Certificate{
		Leaf:        cert,
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}
