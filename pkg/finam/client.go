package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
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

			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.Canceled {
					slog.Debug("finam access token stream exited", "account", c.accountID, "reason", "context cancelled")
					return
				}
			}

			slog.Error("finam acceess token stream exited with error", "account", c.accountID, "msg", err.Error())

			attempts := 0
			sleep := 100 * time.Millisecond

			for {
				attempts++

				stream, err = c.authService.SubscribeJwtRenewal(c.ctx, &auth.SubscribeJwtRenewalRequest{
					Secret: c.secret,
				})

				if err == nil {
					slog.Debug("finam access token stream successfully reconnected", "account", c.accountID, "attempts", attempts)
					break
				}

				slog.Error("cannot reconnect to finam access token stream", "account", c.accountID, "attempt", attempts, "msg", err.Error())
				if attempts > 10 {
					slog.Error("finam access token stream reached the max reconnection attempts", "account", c.accountID, "attempts", attempts)
					c.Close()
					return
				}

				time.Sleep(sleep)
				sleep = sleep * 2
			}
		}
	}()

	return nil
}
