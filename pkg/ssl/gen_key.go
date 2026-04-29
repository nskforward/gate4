package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"

	"github.com/nskforward/gate4/pkg/dir"
)

func GenKey(dst string) (*rsa.PrivateKey, error) {
	f, err := os.Create(dir.Normalize(dst))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, &buf)
	return key, err
}
