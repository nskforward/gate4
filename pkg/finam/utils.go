package finam

import (
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/gate4/pkg/types"
	"github.com/shopspring/decimal"
	protodecimal "google.golang.org/genproto/googleapis/type/decimal"
)

func convertPositions(in []*accounts.Position) []types.Position {
	result := make([]types.Position, 0, len(in))
	for _, pos := range in {
		result = append(result, convertPosition(pos))
	}
	return result
}

func convertPosition(in *accounts.Position) types.Position {
	return types.Position{
		Symbol: in.Symbol,
		Price:  in.AveragePrice.Value,
		Size:   in.Quantity.Value,
		Profit: in.UnrealizedPnl.Value,
	}
}

func convertAsset(in *assets.GetAssetResponse) types.AssetInfo {
	return types.AssetInfo{
		Symbol:   in.Ticker,
		Name:     in.Name,
		Decimals: in.Decimals,
		Currency: in.QuoteCurrency,
		LotSize:  in.LotSize.Value,
		MinStep:  calcMinStep(in.MinStep, in.Decimals),
	}
}

func convertSessions(in []*assets.ScheduleResponse_Sessions) []types.Session {
	result := make([]types.Session, 0, len(in))
	for _, sess := range in {
		result = append(result, convertSession(sess))
	}
	return result
}

func convertSession(in *assets.ScheduleResponse_Sessions) types.Session {
	return types.Session{
		Type:  convertSessionType(in.Type),
		Start: in.Interval.StartTime.Seconds,
		End:   in.Interval.EndTime.Seconds,
	}
}

func convertSessionType(sessionType string) types.SessionType {
	switch sessionType {
	case "CLOSED":
		// no any orders
		return types.SessionClosed
	case "OPENING_AUCTION":
		// only limit orders
		return types.SessionPremarket
	case "EARLY_TRADING":
		// only limit prders
		return types.SessionPremarket
	case "CORE_TRADING":
		// any orders
		return types.SessionMain
	case "CLOSING_AUCTION":
		// only limit orders
		return types.SessionPostmarket
	case "LATE_TRADING":
		// any orders
		return types.SessionPostmarket
	default:
		// no any orders
		return types.SessionUnspecified
	}
}

func convertQuotes(in []*marketdata.Quote) []types.Quote {
	result := make([]types.Quote, 0, len(in))
	for _, q := range in {
		result = append(result, convertQuote(q))
	}
	return result
}

func convertQuote(in *marketdata.Quote) types.Quote {
	return types.Quote{
		Timestamp: in.Timestamp.Seconds,
		Symbol:    in.Symbol,
		Ask: types.QuoteLevel{
			Price: getNullableDecimal(in.Ask),
			Size:  getNullableDecimal(in.AskSize),
		},
		Bid: types.QuoteLevel{
			Price: getNullableDecimal(in.Bid),
			Size:  getNullableDecimal(in.BidSize),
		},
	}
}

func getNullableDecimal(v *protodecimal.Decimal) string {
	if v == nil {
		return "0"
	}
	return v.Value
}

func calcMinStep(minStep int64, decimals int32) string {
	// min_step/(10ˆdecimals)
	minStepDec := decimal.NewFromInt(minStep)
	base := decimal.NewFromInt(10)
	exp := decimal.NewFromInt32(decimals)
	divisor := base.Pow(exp)
	return minStepDec.Div(divisor).StringFixed(decimals)
}

func changedQuote(cache map[string]*types.Quote, q types.Quote) *types.Quote {
	cached, ok := cache[q.Symbol]
	if !ok {
		cache[q.Symbol] = &q
		return &q
	}

	updated := false

	if q.Ask.Price != "0" && q.Ask.Price != cached.Ask.Price {
		cached.Ask.Price = q.Ask.Price
		updated = true
	}

	if q.Bid.Price != "0" && q.Bid.Price != cached.Bid.Price {
		cached.Bid.Price = q.Bid.Price
		updated = true
	}

	if !updated {
		return nil
	}

	return cached
}
