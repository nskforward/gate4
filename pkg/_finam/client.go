package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	ctx        context.Context
	cancel     context.CancelFunc
	account    string
	tokenStore *TokenStore
	errStore   atomic.Pointer[error]
	service    struct {
		marketdata marketdata.MarketDataServiceClient
	}
}

func newClient(ctx context.Context, account, secret string) (*Client, error) {

	conn, err := connect()
	if err != nil {
		return nil, fmt.Errorf("finam client dial error: %w", err)
	}

	clientCtx, cancel := context.WithCancel(ctx)

	client := &Client{
		ctx:        clientCtx,
		cancel:     cancel,
		account:    account,
		tokenStore: NewTokenStore(clientCtx, cancel, account, secret, conn),
	}

	client.service.marketdata = marketdata.NewMarketDataServiceClient(conn)

	err = client.refreshToken()
	if err != nil {
		return nil, fmt.Errorf("finam client auth error: %w", err)
	}

	return client, nil
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) Err() error {
	err := c.errStore.Load()
	if err != nil {
		return *err
	}
	return nil
}

func (c *Client) authCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return c.tokenStore.RequestContext(ctx)
}

func (c *Client) refreshToken() error {
	return c.tokenStore.RefreshToken()
}

func connect() (*grpc.ClientConn, error) {
	return grpc.NewClient("api.finam.ru:443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})),
		grpc.WithIdleTimeout(10*time.Minute),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             time.Minute,
			PermitWithoutStream: true,
		}),
	)
}
