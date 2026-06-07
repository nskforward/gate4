package server

import (
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/ssl"
)

func NewTLSConfig(cfg config.Config) (*tls.Config, error) {

	cert, err := ssl.LoadCertificate(cfg.SSL.Server.Cert)
	if err != nil {
		return nil, fmt.Errorf("cannot load server cert: %w", err)
	}

	slog.Debug("tls server cert", "not after", cert.NotAfter.String())

	key, err := ssl.LoadPrivateKey(cfg.SSL.Server.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot load server key: %w", err)
	}

	tlsCert := tls.Certificate{
		Leaf:        cert,
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}, nil
}
