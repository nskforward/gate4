package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	math "math/rand"
	"net"
	"os"
	"time"
)

func CreateCertificate(isCA bool, subject pkix.Name, expires time.Time, addresses []string) *x509.Certificate {
	dnsList := []string{}
	ipList := []net.IP{}
	for _, addr := range addresses {
		if ip := net.ParseIP(addr); ip != nil {
			ipList = append(ipList, ip)
		} else {
			dnsList = append(dnsList, addr)
		}
	}
	return &x509.Certificate{
		SerialNumber:          big.NewInt(math.Int63()),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              expires,
		IsCA:                  isCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IPAddresses:           ipList,
		DNSNames:              dnsList,
	}
}

func SignCertificate(template, ca *x509.Certificate, key, caKey *rsa.PrivateKey) (*x509.Certificate, error) {
	data, err := x509.CreateCertificate(rand.Reader, template, ca, &key.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("cannot sign cert: %w", err)
	}
	return x509.ParseCertificate(data)
}

func LoadCertificate(path string) (*x509.Certificate, error) {
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
	return x509.ParseCertificate(block.Bytes)
}

func SaveCertificate(cert *x509.Certificate, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(f, &buf)
	return err
}
