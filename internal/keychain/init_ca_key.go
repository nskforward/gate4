package keychain

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nskforward/gate4/pkg/console"
)

func initCAKey(ctx context.Context, path string) (*rsa.PrivateKey, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return generatePrivateKey(ctx, path)
		}
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("ca key cannot be a dir: %s", path)
	}
	return loadPrivateKey(path)
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v, _ := pem.Decode(b)
	if v == nil {
		return nil, fmt.Errorf("cannot decode pem file: %s", path)
	}

	if v.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("pem type must be 'RSA PRIVATE KEY'")
	}

	slog.Info("loaded CA private key", "file", path)

	return x509.ParsePKCS1PrivateKey(v.Bytes)
}

func generatePrivateKey(ctx context.Context, path string) (*rsa.PrivateKey, error) {
	scanner := console.NewScanner()
	defer scanner.Close()

	fmt.Println("WARNING! CA private key not found at", path)

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
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, &buf)
	if err != nil {
		return nil, err
	}

	fmt.Println("success: private key generated at", path)

	return key, err
}
