package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	math "math/rand"
	"os"
	"time"
)

func GenCA(dstKey, dstCert string) error {
	privateKey, err := GenKey(dstKey)
	if err != nil {
		return fmt.Errorf("cannot create private key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(math.Int63()),
		Subject: pkix.Name{
			CommonName: "ROOT CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
	}

	data, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("cannot create cert: %w", err)
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: data,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(dstCert, buf.Bytes(), 0755)
	if err != nil {
		return fmt.Errorf("cannot wite cert: %w", err)
	}

	return nil
}
