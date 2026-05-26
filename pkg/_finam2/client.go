package finam

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/nskforward/gate4/pkg/types"
)

type Client struct {
	ctx        context.Context
	account    string
	conn       *Conn
	tokenStore *TokenStore
}

func NewClient(ctx context.Context, account, secret string) (*Client, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}
	client := &Client{
		ctx:        ctx,
		account:    account,
		conn:       conn,
		tokenStore: NewTokenStore(conn, secret),
	}
	return client, nil
}

func (client *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	reqCtx, cancel := joinContext(client.ctx, ctx)
	defer cancel()

	for {
		err := client.conn.SubscribeQuotes(reqCtx, client.tokenStore, symbol, send)
		if errors.Is(err, ErrPeerAway) {
			slog.Debug("stop finam quote stream", "account", client.account, "symbol", symbol, "reason", err.Error())
			return nil
		}
		slog.Debug("try to reconnect to finam quote stream", "account", client.account, "symbol", symbol, "reason", err.Error())
		time.Sleep(time.Second)
	}
}
