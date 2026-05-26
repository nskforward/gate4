package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Conn struct {
	c          *grpc.ClientConn
	auth       auth.AuthServiceClient
	marketdata marketdata.MarketDataServiceClient
}

func Connect() (*Conn, error) {
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

	return &Conn{
		c:          c,
		auth:       auth.NewAuthServiceClient(c),
		marketdata: marketdata.NewMarketDataServiceClient(c),
	}, nil
}

func (conn *Conn) Close() {
	conn.c.Close()
}

func (conn *Conn) Authorize(ctx context.Context, secret string) (Token, error) {

	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tokenResp, err := conn.auth.Auth(reqCtx, &auth.AuthRequest{
		Secret: secret,
	})
	if err != nil {
		return Token{}, err
	}

	infoResp, err := conn.auth.TokenDetails(reqCtx, &auth.TokenDetailsRequest{
		Token: tokenResp.GetToken(),
	})
	if err != nil {
		return Token{}, err
	}

	return NewToken(ctx, conn.marketdata, tokenResp.GetToken(), infoResp.GetCreatedAt().AsTime(), infoResp.GetExpiresAt().AsTime()), nil
}

func (conn *Conn) SubscribeQuotes(ctx context.Context, tokenStore *TokenStore, symbol string, send func(types.Quote) error) error {
	token, err := tokenStore.GetToken(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return ErrPeerAway
		}
		return fmt.Errorf("cannot get finam auth token: %w", err)
	}

	reqCtx, cancel := token.AutorizedContext(ctx)
	defer cancel()

	stream, err := conn.marketdata.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
		Symbols: []string{symbol},
	})
	if err != nil {
		if ctx.Err() != nil {
			return ErrPeerAway
		}
		return fmt.Errorf("cannot subscribe for finam quote stream: %w", err)
	}

	cache := NewQuoteCache()

	for {
		resp, err := stream.Recv()
		if err != nil {
			if ctx.Err() != nil {
				return ErrPeerAway
			}
			if strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
				_, err = tokenStore.RefreshToken(ctx, token)
				if err != nil {
					slog.Error("cannot refresh finam auth token", "reason", err.Error())
				}
				return ErrUnauthorized
			}
			fmt.Println("token ctx:", token.ctx.Err())
			return err // expect reconnect
		}

		if len(resp.Quote) == 0 {
			continue
		}

		q := resp.Quote[0]

		if !cache.Allow(q) {
			continue
		}

		err = send(types.Quote{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp.Seconds,
			Ask:       cache.Ask,
			Bid:       cache.Bid,
		})
		if err != nil {
			return err
		}
	}
}
