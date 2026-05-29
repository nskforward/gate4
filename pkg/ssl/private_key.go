package ssl

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GeneratePrivateKeyRSA() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

func LoadPrivateKey(path string) (crypto.PrivateKey, error) {
	normalizedPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(normalizedPath)
	if err != nil {
		return nil, err
	}
	return ParsePrivateKey(data)
}

func ParsePrivateKey(data []byte) (crypto.PrivateKey, error) {
	v, _ := pem.Decode(data)
	if v == nil {
		return nil, fmt.Errorf("cannot decode pem file")
	}
	switch v.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(v.Bytes)
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", v.Type)
	}
}

func SavePrivateKeyRSA(key *rsa.PrivateKey, path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

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
