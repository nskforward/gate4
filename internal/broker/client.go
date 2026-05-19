package broker

import (
	"context"

	"github.com/nskforward/gate4/pkg/types"
)

type Client interface {
	GetAccountInfo() types.AccountInfo
	GetPositions(ctx context.Context) ([]types.Position, error)
	GetAsset(ctx context.Context, symbol string) (types.AssetInfo, error)
	GetSchedule(ctx context.Context, symbol string) ([]types.Session, error)
	SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) bool) error
	SubscribeAccountTrades(ctx context.Context, send func(types.AccountTrade) bool) error
}
