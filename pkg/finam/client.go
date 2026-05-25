package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	accountID         string
	secret            string
	tokenStore        *TokenStore
	authService       auth.AuthServiceClient
	markedDataService marketdata.MarketDataServiceClient
	ctx               context.Context
	cancel            context.CancelFunc
}

func newClient(ctx context.Context, accountID, secret string) (*Client, error) {
	conn, err := grpc.NewClient("api.finam.ru:443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})),
		grpc.WithIdleTimeout(5*time.Minute),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Minute,
			Timeout:             30 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("finam client dial error: %w", err)
	}

	clientCtx, cancel := context.WithCancel(ctx)

	client := &Client{
		ctx:               clientCtx,
		cancel:            cancel,
		accountID:         accountID,
		secret:            secret,
		authService:       auth.NewAuthServiceClient(conn),
		markedDataService: marketdata.NewMarketDataServiceClient(conn),
		tokenStore:        NewTokenStore(clientCtx),
	}

	err = client.subscribeTokenRefresh()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) subscribeTokenRefresh() error {
	stream, err := c.authService.SubscribeJwtRenewal(c.ctx, &auth.SubscribeJwtRenewalRequest{
		Secret: c.secret,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err == nil {
				c.tokenStore.set(resp.GetToken())
				slog.Debug("refreshed finam token", "account", c.accountID)
				continue
			}

			if tools.IsGRPCCancelled(err) {
				slog.Debug("finam access token stream aborted", "account", c.accountID, "reason", "context cancelled")
				return
			}

			slog.Error("finam acceess token stream aborted", "account", c.accountID, "reason", err.Error())

			retry := NewRetry(func() (grpc.ServerStreamingClient[auth.SubscribeJwtRenewalResponse], error) {
				return c.authService.SubscribeJwtRenewal(c.ctx, &auth.SubscribeJwtRenewalRequest{
					Secret: c.secret,
				})
			})

			retry.OnSuccess = func(attempt int) {
				slog.Debug("finam access token stream successfully connected", "account", c.accountID, "attempt", attempt)
			}

			retry.OnFailure = func(err error, attempt int) {
				slog.Error("finam access token stream connection attempt failed", "account", c.accountID, "attempt", attempt, "msg", err.Error())
			}

			stream, err = retry.Do(c.ctx)
			if err != nil {
				slog.Error("access token stream reached the max reconnection attemps", "account", c.accountID, "last_error", err)
				c.Close()
				return
			}
		}
	}()

	return nil
}
