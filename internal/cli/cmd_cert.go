package cli

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/pkg/ssl"
)

func (c *Client) cmdCert(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("requres sub command")
	}

	command := args[0]
	args = args[1:]

	switch command {
	case "init-ca":
		return c.cmdCertInitCA(ctx, args)
	case "issue":
		return c.cmdCertIssue(ctx, args)
	case "list-active":
		return c.cmdCertListActive(ctx, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (c *Client) cmdCertInitCA(_ context.Context, _ []string) error {
	commonName := Ask("common name")
	certPath := AskPath("dest CA cert")
	keyPath := AskPath("dest CA key")
	err := ssl.GenCA(keyPath, certPath, pkix.Name{
		CommonName: commonName,
	})
	if err != nil {
		return fmt.Errorf("cannot generate CA: %w", err)
	}
	fmt.Println("success")
	return nil
}

func (c *Client) cmdCertIssue(_ context.Context, _ []string) error {
	commonName := Ask("common name")
	certPath := AskPath("dest cert")
	keyPath := AskPath("dest key")
	caCertPath := AskPath("src CA cert")
	caKeyPath := AskPath("src CA key")
	err := ssl.GenCert(caKeyPath, caCertPath, keyPath, certPath, pkix.Name{
		CommonName: commonName,
	})
	if err != nil {
		return fmt.Errorf("cannot generate ssl cert: %w", err)
	}
	fmt.Println("success")
	return nil
}

func (c *Client) cmdCertListActive(ctx context.Context, args []string) error {
	return errors.New("not implemented")
}
