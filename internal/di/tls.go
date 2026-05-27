package di

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/console"
)

func initTLSConfig(cfg *config.Config) *tls.Config {
	loadedServerCert, err := tls.LoadX509KeyPair(cfg.Server.SSL.Cert, cfg.Server.SSL.Key)
	if err != nil {
		console.LogFatal("cannot load server cert", err)
	}

	loadedCACert, err := os.ReadFile(cfg.CA.Cert)
	if err != nil {
		console.LogFatal("cannot load CA cert", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(loadedCACert) {
		console.LogFatal("cannot add CA cert to the pool", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{loadedServerCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}
	return tlsConfig
}
