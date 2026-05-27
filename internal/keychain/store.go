package keychain

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"

	"github.com/nskforward/gate4/pkg/ssl"
)

type Store struct {
	caKey  crypto.PrivateKey
	caCert *x509.Certificate
}

func NewStore(caKey crypto.PrivateKey, caCert *x509.Certificate) *Store {
	return &Store{
		caKey:  caKey,
		caCert: caCert,
	}
}

func (store *Store) Generate(commonName string, key crypto.PrivateKey) (*x509.Certificate, error) {
	template := ssl.CreateCertificate(false, time.Now().AddDate(2, 0, 0), nil, pkix.Name{
		CommonName:   commonName,
		Organization: []string{"Gate4 LLC"},
		Country:      []string{"RU"},
	})
	return ssl.SignCertificate(template, store.caCert, key.(crypto.Signer), store.caKey)
}
