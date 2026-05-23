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
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/pkg/console"
)

func initServerPrivateKey(ctx context.Context, path string) (*rsa.PrivateKey, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return generatePrivateKey(ctx, path)
		}
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("server key cannot be a dir at %s", path)
	}
	return loadPrivateKey(path)
}

func initServerCert(ctx context.Context, path string, key *rsa.PrivateKey) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = generateCert(ctx, path, key)
			return err
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("ca cert cannot be a dir: %s", path)
	}
	return nil
}

func generateCert(ctx context.Context, path string, key *rsa.PrivateKey) (*x509.Certificate, error) {
	scanner := console.NewScanner()
	defer scanner.Close()

	fmt.Println("WARNING! server cert not found at", path)

	yes, err := scanner.ScanBool(ctx, "generate?", nil, nil)
	if err != nil {
		return nil, err
	}

	if !yes {
		return nil, fmt.Errorf("aborted by user")
	}

	dnsNames := make([]string, 0, 1)
	ipAddresses := make([]net.IP, 0, 1)

	commonName, err := scanner.Scan(ctx, "common name", "", nil)
	if err != nil {
		return nil, err
	}

	dnsAddr, err := scanner.Scan(ctx, "dns address", "", nil)
	if err != nil {
		return nil, err
	}
	if dnsAddr != "" {
		dnsNames = append(dnsNames, dnsAddr)
	}

	ipAddr, err := scanner.Scan(ctx, "ip address", "", nil)
	if err != nil {
		return nil, err
	}

	if ipAddr != "" {
		ipAddresses = append(ipAddresses, net.ParseIP(ipAddr))
	}

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
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
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IPAddresses:           ipAddresses,
		DNSNames:              dnsNames,
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
