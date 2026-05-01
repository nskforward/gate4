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

	"github.com/nskforward/gate4/pkg/dir"
)

func GenCert(srcCAKey, srcCACert, dstKey, dstCert string, subject pkix.Name) error {

	caKey, err := LoadKey(srcCAKey)
	if err != nil {
		return fmt.Errorf("cannot load CA key: %w", err)
	}

	caCert, err := LoadCert(srcCACert)
	if err != nil {
		return fmt.Errorf("cannot load CA cert: %w", err)
	}

	privateKey, err := GenKey(dir.Normalize(dstKey))
	if err != nil {
		return fmt.Errorf("cannot create private key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(math.Int63()),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
	}

	data, err := x509.CreateCertificate(rand.Reader, template, caCert, &privateKey.PublicKey, caKey)
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

	err = os.WriteFile(dir.Normalize(dstCert), buf.Bytes(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot wite cert: %w", err)
	}

	return nil
}
