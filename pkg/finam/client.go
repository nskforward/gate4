package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	ctx        context.Context
	cancel     context.CancelFunc
	accountID  string
	secret     string
	tokenStore *TokenStore
	errStore   atomic.Pointer[error]
	service    struct {
		auth       auth.AuthServiceClient
		markeddata marketdata.MarketDataServiceClient
	}
}

func newClient(ctx context.Context, accountID, secret string) (*Client, error) {

	conn, err := connect()
	if err != nil {
		return nil, fmt.Errorf("finam client dial error: %w", err)
	}

	clientCtx, cancel := context.WithCancel(ctx)

	client := &Client{
		ctx:        clientCtx,
		cancel:     cancel,
		accountID:  accountID,
		secret:     secret,
		tokenStore: NewTokenStore(clientCtx),
	}

	client.service.auth = auth.NewAuthServiceClient(conn)
	client.service.markeddata = marketdata.NewMarketDataServiceClient(conn)

	err = client.refreshToken()
	if err != nil {
		return nil, fmt.Errorf("finam client auth error: %w", err)
	}

	go client.watchToken()

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
