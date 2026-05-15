package finam

import (
	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
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

func convertAccountTrade(in *v1.AccountTrade) types.AccountTrade {
	size := in.Size.Value
	if in.Side == v1.Side_SIDE_SELL {
		size = "-" + size
	}
	return types.AccountTrade{
		Timestamp: in.Timestamp.Seconds,
		AccountID: in.AccountId,
		Symbol:    in.Symbol,
		Price:     in.Price.Value,
		Size:      size,
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
