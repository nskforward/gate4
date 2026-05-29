package ssl

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	math "math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/pkg/common"
)

func CreateCertificate(isCA bool, expires time.Time, addresses []string, subject pkix.Name) *x509.Certificate {
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

func SignCertificate(template, ca *x509.Certificate, key crypto.Signer, caKey crypto.PrivateKey) (*x509.Certificate, error) {
	data, err := x509.CreateCertificate(rand.Reader, template, ca, key.Public(), caKey)
	if err != nil {
		return nil, fmt.Errorf("cannot sign cert: %w", err)
	}
	return x509.ParseCertificate(data)
}

func LoadCertificate(path string) (*x509.Certificate, error) {
	data, err := os.ReadFile(common.Path(path))
	if err != nil {
		return nil, err
	}
	return ParseCertificate(data)
}

func ParseCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("certificate must be in PEM format")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("pem file is not certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func SaveCertificate(cert *x509.Certificate, path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	r, err := MarshalCert(cert)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

func MarshalCert(cert *x509.Certificate) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	return &buf, err
}
