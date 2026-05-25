package finam

import (
	"context"
	"log/slog"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	reqCtx, cancel := c.tokenStore.RequestContext(ctx)
	defer cancel()

	stream, err := c.markedDataService.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
		Symbols: []string{symbol},
	})
	if err != nil {
		return err
	}

	cache := NewQuoteCache()

	for {
		resp, err := stream.Recv()

		if err == nil {
			if len(resp.Quote) == 0 {
				continue
			}
			q := resp.Quote[0]

			if !cache.Allow(q) {
				continue
			}

			err := send(types.Quote{
				Symbol:    q.Symbol,
				Timestamp: q.Timestamp.Seconds,
				Ask:       cache.Ask,
				Bid:       cache.Bid,
			})
			if err != nil {
				return err
			}
			continue
		}

		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.Canceled {
				slog.Debug("finam quote stream exited", "account", c.accountID, "symbol", symbol, "reason", "context cancelled")
				return nil
			}
		}

		slog.Error("finam quote stream exited with error", "account", c.accountID, "symbol", symbol, "msg", err.Error())

		attempts := 0
		sleep := 100 * time.Millisecond

		for {
			attempts++

			cancel() // close prev reqCtx

			reqCtx, cancel = c.tokenStore.RequestContext(ctx)
			stream, err = c.markedDataService.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
				Symbols: []string{symbol},
			})

			if err == nil {
				slog.Debug("finam quote stream successfully reconnected", "account", c.accountID, "symbol", symbol, "attempts", attempts)
				break
			}

			slog.Error("cannot reconnect to finam quote stream", "account", c.accountID, "symbol", symbol, "attempt", attempts, "msg", err.Error())

			if attempts > 10 {
				slog.Error("finam quote stream reached the max reconnection attempts", "account", c.accountID, "symbol", symbol, "attempts", attempts)
				return err
			}

			time.Sleep(sleep)
			sleep = sleep * 2
		}
	}
}
