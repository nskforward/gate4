package types

import (
	"context"
)

type Client interface {
	GetLastQuote(ctx context.Context, symbol string) (Quote, error)
	GetAccountInfo() AccountInfo
	GetPositions(ctx context.Context) ([]Position, error)
	GetAsset(ctx context.Context, symbol string) (AssetInfo, error)
	GetSchedule(ctx context.Context, symbol string) ([]Session, error)
	SubscribeQuotes(ctx context.Context, symbol string, send func(Quote) bool) error
	SubscribeAccountTrades(ctx context.Context, send func(AccountTrade) bool) error
}
