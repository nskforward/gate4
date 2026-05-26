package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Conn struct {
	grpcConn *grpc.ClientConn
	creds    *Creds
	auth     auth.AuthServiceClient
	token    *Token
	tokenMx  sync.Mutex
}

func Connect(creds *Creds) (*Conn, error) {
	c, err := grpc.NewClient("api.finam.ru:443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})),
		grpc.WithIdleTimeout(10*time.Minute),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             time.Minute,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		grpcConn: c,
		creds:    creds,
		auth:     auth.NewAuthServiceClient(c),
	}

	_, err = conn.GetToken()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (conn *Conn) Close() error {
	return conn.grpcConn.Close()
}

func (conn *Conn) GetToken() (*Token, error) {
	conn.tokenMx.Lock()
	defer conn.tokenMx.Unlock()

	if conn.token != nil && time.Until(conn.token.Expires) > time.Minute {
		return conn.token, nil
	}

	t, err := conn.createToken()
	if err != nil {
		return nil, fmt.Errorf("finam auth failed for account %s: %w", conn.creds.AccountID, err)
	}

	conn.token = t

	return t, nil
}

func (conn *Conn) createToken() (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokenResp, err := conn.auth.Auth(ctx, &auth.AuthRequest{
		Secret: conn.creds.Secret,
	})
	if err != nil {
		return nil, err
	}

	infoResp, err := conn.auth.TokenDetails(ctx, &auth.TokenDetailsRequest{
		Token: tokenResp.GetToken(),
	})
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:   tokenResp.GetToken(),
		Created: time.Unix(infoResp.GetCreatedAt().Seconds, 0),
		Expires: time.Unix(infoResp.GetExpiresAt().Seconds, 0),
	}, nil
}
