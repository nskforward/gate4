package keychain

import (
	"context"
	"crypto"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"time"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/ssl"
	"github.com/nskforward/gate4/pkg/tools"
)

func GenWizard(ctx context.Context, cfg *config.Config) error {
	if !tools.FileExists(cfg.CA.Key) || !tools.FileExists(cfg.CA.Cert) {
		caKey, err := generateCAKey(ctx, cfg.CA.Key)
		if err != nil {
			return err
		}
		if err := generateCACert(ctx, cfg.CA.Cert, caKey); err != nil {
			return err
		}
	}

	if !tools.FileExists(cfg.Server.SSL.Key) || !tools.FileExists(cfg.Server.SSL.Cert) {
		key, err := generateServerKey(ctx, cfg.Server.SSL.Key)
		if err != nil {
			return err
		}
		if err := generateServerCert(ctx, cfg.Server.SSL.Cert, key); err != nil {
			return err
		}
	}

	return nil
}

func generateCAKey(ctx context.Context, path string) (crypto.PrivateKey, error) {
	scanner := console.NewScanner()
	defer scanner.Close()
	fmt.Print("checking CA ssl key... ")
	if !tools.FileExists(path) {
		printNotFound()
		fmt.Println("WARNING! A new CA ssl private key will be generated")
		if !scanner.Confirm(ctx, "continue?") {
			return nil, errors.New("aborted")
		}
		caKey, err := ssl.GeneratePrivateKeyRSA()
		if err != nil {
			return nil, err
		}
		err = ssl.SavePrivateKeyRSA(caKey, path)
		return caKey, err
	}
	printOk()
	return ssl.LoadPrivateKey(path)
}

func generateCACert(ctx context.Context, path string, key crypto.PrivateKey) error {
	scanner := console.NewScanner()
	defer scanner.Close()
	fmt.Print("checking CA ssl cert... ")
	if !tools.FileExists(path) {
		printNotFound()
		fmt.Println("WARNING! A new CA ssl cert key will be generated")
		if !scanner.Confirm(ctx, "continue?") {
			return errors.New("aborted")
		}
		template := ssl.CreateCertificate(true, time.Now().AddDate(10, 0, 0), nil, pkix.Name{
			CommonName: "Gate4 Root CA 1",
		})
		cert, err := ssl.SignCertificate(template, template, key.(crypto.Signer), key)
		if err != nil {
			return err
		}
		return ssl.SaveCertificate(cert, path)
	}
	printOk()
	return nil
}

func generateServerKey(ctx context.Context, path string) (crypto.PrivateKey, error) {
	scanner := console.NewScanner()
	defer scanner.Close()
	fmt.Print("checking server ssl key... ")
	if !tools.FileExists(path) {
		printNotFound()
		fmt.Println("WARNING! A new server ssl private key will be generated")
		if !scanner.Confirm(ctx, "continue?") {
			return nil, errors.New("aborted")
		}
		key, err := ssl.GeneratePrivateKeyRSA()
		if err != nil {
			return nil, err
		}
		err = ssl.SavePrivateKeyRSA(key, path)
		return key, err
	}
	printOk()
	return ssl.LoadPrivateKey(path)
}

func generateServerCert(ctx context.Context, path string, key crypto.PrivateKey) error {
	scanner := console.NewScanner()
	defer scanner.Close()
	fmt.Print("checking server ssl cert... ")
	if !tools.FileExists(path) {
		printNotFound()
		fmt.Println("WARNING! A new server ssl cert key will be generated")
		if !scanner.Confirm(ctx, "continue?") {
			return errors.New("aborted")
		}
		template := ssl.CreateCertificate(true, time.Now().AddDate(10, 0, 0), nil, pkix.Name{
			CommonName: "Gate4 Server",
		})
		cert, err := ssl.SignCertificate(template, template, key.(crypto.Signer), key)
		if err != nil {
			return err
		}
		return ssl.SaveCertificate(cert, path)
	}
	printOk()
	return nil
}

func printNotFound() {
	fmt.Println(console.FormatText("not found", console.Red))
}

func printOk() {
	fmt.Println("ok")
}
