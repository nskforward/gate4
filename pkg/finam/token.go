package finam

import (
	"context"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"google.golang.org/grpc/metadata"
)

type Token struct {
	ctx        context.Context
	cancel     context.CancelFunc
	value      string
	created    time.Time
	expires    time.Time
	marketdata marketdata.MarketDataServiceClient
}

func NewToken(ctx context.Context, marketdata marketdata.MarketDataServiceClient, value string, created, expires time.Time) Token {
	ctx, cancel := context.WithCancel(ctx)
	return Token{
		ctx:        ctx,
		cancel:     cancel,
		value:      value,
		created:    created,
		expires:    expires,
		marketdata: marketdata,
	}
}

func (t Token) Close() {
	t.cancel()
}

func (t Token) BeforeExpiration() time.Duration {
	return time.Until(t.expires)
}

func (t Token) AutorizedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	authCtx := metadata.AppendToOutgoingContext(ctx, "Authorization", t.value)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-t.ctx.Done():
			cancel()
			return
		case <-time.After(time.Until(t.expires) - time.Second):
			cancel()
			return
		}
	}()
	return authCtx, cancel
}
