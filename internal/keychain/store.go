package keychain

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"time"

	"github.com/nskforward/gate4/pkg/ssl"
)

type Store struct {
	caKey  *rsa.PrivateKey
	caCert *x509.Certificate
}

func NewStore(caKey, caCert string) (*Store, error) {
	return &Store{}, errors.New("not implemented")
}

func (store *Store) Generate(commonName string, key *rsa.PrivateKey) (*x509.Certificate, error) {
	template := ssl.CreateCertificate(false, time.Now().AddDate(2, 0, 0), nil, pkix.Name{
		CommonName:   commonName,
		Organization: []string{"Gate4 LLC"},
		Country:      []string{"RU"},
	})
	return ssl.SignCertificate(template, store.caCert, key, store.caKey)
}
