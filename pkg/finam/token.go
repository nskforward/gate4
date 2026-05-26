package finam

import (
	"context"
	"log/slog"
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

func NewToken(marketdata marketdata.MarketDataServiceClient, value string, created, expires time.Time) Token {
	ctx, cancel := context.WithCancel(context.Background())
	t := Token{
		ctx:        ctx,
		cancel:     cancel,
		value:      value,
		created:    created,
		expires:    expires,
		marketdata: marketdata,
	}
	go t.watch()
	slog.Debug("finam auth token created", "expires", time.Unix(expires.Unix(), 0).Format("2006-01-02 15:04"))
	return t
}

func (t Token) watch() {
	defer t.cancel()

	select {
	case <-t.ctx.Done():
		return
	case <-time.After(time.Until(t.expires) - time.Minute):
		return
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
		}
	}()
	return authCtx, cancel
}
