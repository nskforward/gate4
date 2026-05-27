package di

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
)

func initTLSConfig(serverCert *x509.Certificate, serverKey crypto.PrivateKey, caCerts ...*x509.Certificate) *tls.Config {
	tlsCert := tls.Certificate{
		Certificate: [][]byte{serverCert.Raw},
		PrivateKey:  serverKey,
	}
	tlsCert.Leaf = serverCert

	caCertPool := x509.NewCertPool()
	for _, ca := range caCerts {
		caCertPool.AddCert(ca)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}
}
