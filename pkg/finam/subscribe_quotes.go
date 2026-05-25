package finam

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/pkg/types"
)

func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	return errors.New("not implemented")
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
