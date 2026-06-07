package handler

import (
	"context"
	"crypto/x509/pkix"
	"fmt"
	"time"

	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/ssl"
)

func CreateCert(ctx context.Context, args []string) error {
	scanner := input.NewScanner()
	defer scanner.Close()

	caKeyPath, err := scanner.Scan(ctx, "ca key path", "", nil)
	if err != nil {
		return err
	}

	caKey, err := ssl.LoadPrivateKey(caKeyPath)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	caCerPath, err := scanner.Scan(ctx, "ca cert path", "", nil)
	if err != nil {
		return err
	}

	caCert, err := ssl.LoadCertificate(caCerPath)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	privateKeyPath, err := scanner.Scan(ctx, "private key path", "", nil)
	if err != nil {
		return err
	}

	key, err := ssl.GeneratePrivateKeyRSA()
	if err != nil {
		return err
	}

	err = ssl.SavePrivateKeyRSA(key, privateKeyPath)
	if err != nil {
		return err
	}

	commonName, err := scanner.Scan(ctx, "common name", "", nil)
	if err != nil {
		return err
	}

	certPath, err := scanner.Scan(ctx, "cert path", "", nil)
	if err != nil {
		return err
	}

	tpl := ssl.CreateCertificate(false, time.Now().AddDate(2, 0, 0), []string{"localhost"}, pkix.Name{
		CommonName: commonName,
	})

	cert, err := ssl.SignCertificate(tpl, caCert, key, caKey)
	if err != nil {
		return err
	}

	err = ssl.SaveCertificate(cert, certPath)
	if err != nil {
		return err
	}

	return nil
}
