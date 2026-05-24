package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	accountID         string
	authService       auth.AuthServiceClient
	markedDataService marketdata.MarketDataServiceClient
	token             string
}

func newClient(accountID, secret string) (*Client, error) {
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

	authService := auth.NewAuthServiceClient(conn)
	markedDataService := marketdata.NewMarketDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := authService.Auth(ctx, &auth.AuthRequest{Secret: secret})
	if err != nil {
		return nil, fmt.Errorf("finam client auth error: %w", err)
	}

	return &Client{
		accountID:         accountID,
		authService:       authService,
		markedDataService: markedDataService,
		token:             resp.GetToken(),
	}, nil
}

func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	reqCtx := metadata.AppendToOutgoingContext(ctx, "Authorization", c.token)

	stream, err := c.markedDataService.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
		Symbols: []string{symbol},
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.Canceled {
					return nil
				}
			}
			return err
		}
		for _, q := range resp.Quote {
			ask := "0"
			bid := "0"

			if q.Ask != nil {
				ask = q.Ask.Value
			}

			if q.Bid != nil {
				bid = q.Bid.Value
			}

			if err := send(types.Quote{
				Symbol:    q.Symbol,
				Timestamp: q.Timestamp.Seconds,
				Ask:       ask,
				Bid:       bid,
			}); err != nil {
				return err
			}
		}
	}
}
