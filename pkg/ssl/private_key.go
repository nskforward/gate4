package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v, _ := pem.Decode(b)
	if v == nil {
		return nil, fmt.Errorf("cannot decode pem file: %s", path)
	}
	if v.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("pem type must be 'RSA PRIVATE KEY'")
	}
	return x509.ParsePKCS1PrivateKey(v.Bytes)
}

func SavePrivateKey(key *rsa.PrivateKey, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer

	err = pem.Encode(&buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, &buf)
	return err
}
