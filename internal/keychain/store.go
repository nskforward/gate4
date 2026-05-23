package keychain

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	math "math/rand"
	"time"

	"github.com/nskforward/gate4/internal/config"
)

type Store struct {
	key  *rsa.PrivateKey
	cert *x509.Certificate
}

func NewStore(ctx context.Context, cfg config.Config) (*Store, error) {
	key, err := initCAKey(ctx, cfg.Admin.SSL.CAKey)
	if err != nil {
		return nil, fmt.Errorf("cannot init CA key: %w", err)
	}

	cert, err := initCACert(ctx, cfg.Admin.SSL.CACert, key)
	if err != nil {
		return nil, fmt.Errorf("cannot init CA cert: %w", err)
	}

	return &Store{
		key:  key,
		cert: cert,
	}, nil
}

func (store *Store) GenCert(commonName string, privateKey []byte) ([]byte, error) {
	v, _ := pem.Decode(privateKey)
	if v == nil {
		return nil, fmt.Errorf("cannot decode private key")
	}

	if v.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("private key pem type must be 'RSA PRIVATE KEY'")
	}

	key, err := x509.ParsePKCS1PrivateKey(v.Bytes)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(math.Int63()),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Gate4 LLC"},
			Country:      []string{"RU"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
	}

	data, err := x509.CreateCertificate(rand.Reader, template, store.cert, &key.PublicKey, store.key)
	if err != nil {
		return nil, fmt.Errorf("cannot generate cert: %w", err)
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
