package broker

import (
	"context"

	"github.com/nskforward/gate4/pkg/types"
)

type Client interface {
	GetAccountID() string
	GetBrokerID() string
	GetAccount(ctx context.Context) (*types.AccountInfo, error)
	GetAsset(ctx context.Context, symbol string) (types.AssetInfo, error)
	GetSchedule(ctx context.Context, symbol string) ([]types.Session, error)
	SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) bool) error
	SubscribeAccountTrades(ctx context.Context, send func(types.AccountTrade) bool) error
}
