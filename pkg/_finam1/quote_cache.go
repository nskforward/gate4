package finam

import "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"

type QuoteCache struct {
	Ask string
	Bid string
}

func NewQuoteCache() *QuoteCache {
	return &QuoteCache{
		Ask: "0",
		Bid: "0",
	}
}

func (cache *QuoteCache) Allow(q *marketdata.Quote) bool {
	ask := "0"
	bid := "0"

	if q.Ask != nil {
		ask = q.Ask.Value
	}

	if q.Bid != nil {
		bid = q.Bid.Value
	}

	if ask == "0" && bid == "0" {
		return false
	}

	if ask == cache.Ask && bid == cache.Bid {
		return false
	}

	if ask == "0" {
		ask = cache.Ask
	} else {
		cache.Ask = ask
	}

	if bid == "0" {
		bid = cache.Bid
	} else {
		cache.Bid = bid
	}

	if ask == "0" || bid == "0" {
		return false
	}

	return true
}
