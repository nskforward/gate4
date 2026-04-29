package ssl

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadKey(filename string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	v, _ := pem.Decode(b)
	if v == nil {
		return nil, fmt.Errorf("cannot decode pem file: %s", filename)
	}

	if v.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("pem type must be 'RSA PRIVATE KEY'")
	}

	return x509.ParsePKCS1PrivateKey(v.Bytes)
}

func LoadCert(filename string) (*x509.Certificate, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("certificate must be in PEM format: %s", filename)
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("pem file is not certificate: %s", filename)
	}
	return x509.ParseCertificate(block.Bytes)
}
