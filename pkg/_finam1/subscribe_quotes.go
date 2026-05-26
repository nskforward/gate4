package finam

import (
	"context"
	"log/slog"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/tools"
	"github.com/nskforward/gate4/pkg/types"
)

func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	reqCtx, cancel := c.authCtx(ctx)
	defer cancel()

	stream, err := c.service.marketdata.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
		Symbols: []string{symbol},
	})

	if err != nil {
		return err
	}

	cache := NewQuoteCache()

	for {
		resp, err := stream.Recv()

		if err != nil {
			if tools.IsGRPCCancelled(err) {
				slog.Debug("finam quote stream aborted", "account", c.account, "symbol", symbol, "reason", "context cancelled")
				return c.Err()
			}

			slog.Error("finam quote stream aborted", "account", c.account, "symbol", symbol, "msg", err.Error())

			if strings.

			var retry retries.Retry
			for attempt := range retry.Range(ctx) {
				cancel()
				reqCtx, cancel = c.authCtx(ctx)
				stream, err = c.service.marketdata.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
					Symbols: []string{symbol},
				})
				if err != nil {
					slog.Error("finam quote stream connection attempt failed", "account", c.account, "symbol", symbol, "attempt", attempt, "msg", err.Error())
					continue
				}
				break
			}

			if retry.Err() != nil {
				return err
			}

			slog.Debug("finam quote stream successfully reconnected", "account", c.account, "symbol", symbol)

			continue
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

/*
func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	var reqCtx context.Context
	var cancel context.CancelFunc
	cache := NewQuoteCache()

	for {
		retry := NewRetry(func() (grpc.ServerStreamingClient[marketdata.SubscribeQuoteResponse], error) {
			if cancel != nil {
				cancel()
			}
			reqCtx, cancel = c.tokenStore.RequestContext(ctx)
			return c.markedDataService.SubscribeQuote(reqCtx, &marketdata.SubscribeQuoteRequest{
				Symbols: []string{symbol},
			})
		})
		retry.OnSuccess = func(attempt int) {
			slog.Debug("finam quote stream successfully connected", "account", c.accountID, "symbol", symbol, "attempt", attempt)
		}
		retry.OnFailure = func(err error, attempt int) {
			slog.Error("finam quote stream connection attempt failed", "account", c.accountID, "symbol", symbol, "attempt", attempt, "msg", err.Error())
		}

		stream, err := retry.Do(ctx)
		if err != nil {
			return err
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				if tools.IsGRPCCancelled(err) {
					slog.Debug("finam quote stream aborted", "account", c.accountID, "symbol", symbol, "reason", "context cancelled")
					return c.Err()
				}
				slog.Error("finam quote stream aborted", "account", c.accountID, "symbol", symbol, "msg", err.Error())
				break
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
}
*/
