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
	"log/slog"
	"math/big"
	math "math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/pkg/console"
)

func initCACert(ctx context.Context, path string, key *rsa.PrivateKey) (*x509.Certificate, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return generateCertCA(ctx, path, key)
		}
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("ca cert cannot be a dir: %s", path)
	}
	return loadCert(path)
}

func loadCert(path string) (*x509.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("certificate must be in PEM format: %s", path)
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("pem file is not certificate: %s", path)
	}
	slog.Info("loaded CA cert", "file", path)
	return x509.ParseCertificate(block.Bytes)
}

func generateCertCA(ctx context.Context, path string, key *rsa.PrivateKey) (*x509.Certificate, error) {
	scanner := console.NewScanner()
	defer scanner.Close()

	fmt.Println("WARNING! CA cert not found at", path)

	yes, err := scanner.ScanBool(ctx, "generate?", nil, nil)
	if err != nil {
		return nil, err
	}

	if !yes {
		return nil, fmt.Errorf("aborted by user")
	}

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(math.Int63()),
		Subject: pkix.Name{
			CommonName:   "Gate4 Root CA 1",
			Organization: []string{"Gate4 LLC"},
			Country:      []string{"RU"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IPAddresses:           []net.IP{},
		DNSNames:              []string{},
	}

	data, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("cannot create cert: %w", err)
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("cannot parse cert: %w", err)
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(path, buf.Bytes(), 0755)

	fmt.Println("success: cert generated at", path)

	return cert, err
}
