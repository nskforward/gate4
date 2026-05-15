package finam

import (
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/types"
)

type QuoteCache struct {
	quotes map[string]*types.Quote
}

func NewQuoteCache() *QuoteCache {
	return &QuoteCache{
		quotes: make(map[string]*types.Quote),
	}
}

func (cache *QuoteCache) Get(in *marketdata.Quote) *types.Quote {
	if in == nil {
		return nil
	}

	cached, ok := cache.quotes[in.Symbol]
	if !ok {
		result := &types.Quote{
			Timestamp: in.Timestamp.Seconds,
			Symbol:    in.Symbol,
			Ask:       in.Ask.Value,
			Bid:       in.Bid.Value,
		}
		cache.quotes[in.Symbol] = result
		return result
	}

	updated := false

	if in.Ask != nil && in.Ask.Value != cached.Ask {
		cached.Ask = in.Ask.Value
		updated = true
	}

	if in.Bid != nil && in.Bid.Value != cached.Bid {
		cached.Bid = in.Bid.Value
		updated = true
	}

	if updated {
		cached.Timestamp = in.Timestamp.Seconds
		return cached
	}

	return nil
}
