package keychain

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/nskforward/gate4/internal/config"
)

type Store struct {
	key  *rsa.PrivateKey
	cert *x509.Certificate
}

func NewStore(ctx context.Context, cfg config.Config) (*Store, error) {
	key, err := initCAKey(ctx, cfg.Admin.SSL.CAKey)
	if err != nil {
		return nil, fmt.Errorf("cannot init CA key: %w", err)
	}

	cert, err := initCACert(ctx, cfg.Admin.SSL.CACert, key)
	if err != nil {
		return nil, fmt.Errorf("cannot init CA cert: %w", err)
	}

	return &Store{
		key:  key,
		cert: cert,
	}, nil
}
