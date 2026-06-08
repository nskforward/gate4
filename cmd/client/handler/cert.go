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

	isCA, err := scanner.ScanBool(ctx, "is CA", nil, nil)
	if err != nil {
		return err
	}

	if isCA {
		return generateCA(ctx, scanner)
	}

	return generateRegular(ctx, scanner)
}

func generateCA(ctx context.Context, scanner *input.Scanner) error {

	fmt.Println("please specify filename path")

	keyPath, err := scanner.Scan(ctx, "ca key", "", nil)
	if err != nil {
		return err
	}

	key, err := ssl.GeneratePrivateKeyRSA()
	if err != nil {
		return err
	}

	err = ssl.SavePrivateKeyRSA(key, keyPath)
	if err != nil {
		return err
	}

	fmt.Println("+ private key saved")

	certPath, err := scanner.Scan(ctx, "ca cert", "", nil)
	if err != nil {
		return err
	}

	commonName, err := scanner.Scan(ctx, "common name", "", nil)
	if err != nil {
		return err
	}

	tpl := ssl.CreateTemplate(true, time.Now().AddDate(10, 0, 0), nil, pkix.Name{
		CommonName:   commonName,
		Organization: []string{"GATE4 EXCHANGE"},
	})

	cert, err := ssl.CreateCertificate(tpl, tpl, key, key)
	if err != nil {
		return err
	}

	err = ssl.SaveCertificate(cert, certPath)
	if err != nil {
		return err
	}

	fmt.Println("+ certificate saved")

	return nil
}

func generateRegular(ctx context.Context, scanner *input.Scanner) error {

	fmt.Println("please specify the path to the existing CA key and cert files")

	caKeyPath, err := scanner.Scan(ctx, "ca key", "", nil)
	if err != nil {
		return err
	}

	caKey, err := ssl.LoadPrivateKey(caKeyPath)
	if err != nil {
		return err
	}

	fmt.Println("+ ca key loaded")

	caCerPath, err := scanner.Scan(ctx, "ca cert", "", nil)
	if err != nil {
		return err
	}

	caCert, err := ssl.LoadCertificate(caCerPath)
	if err != nil {
		return err
	}

	fmt.Println("+ ca cert loaded")

	fmt.Println("please specify the path where server key and cert files will be saved")

	privateKeyPath, err := scanner.Scan(ctx, "private key", "", nil)
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

	fmt.Println("+ private key saved")

	certPath, err := scanner.Scan(ctx, "cert", "", nil)
	if err != nil {
		return err
	}

	host, err := scanner.Scan(ctx, "host", "", nil)
	if err != nil {
		return err
	}

	tpl := ssl.CreateTemplate(false, time.Now().AddDate(2, 0, 0), []string{host}, pkix.Name{
		CommonName: host,
	})

	cert, err := ssl.CreateCertificate(tpl, caCert, key, caKey)
	if err != nil {
		return err
	}

	err = ssl.SaveCertificate(cert, certPath)
	if err != nil {
		return err
	}

	fmt.Println("+ cert saved")

	return nil
}
