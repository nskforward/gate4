package broker

import (
	"context"

	"github.com/nskforward/gate4/pkg/types"
)

type Client interface {
	GetAccount(ctx context.Context) (*types.AccountInfo, error)
	GetAsset(ctx context.Context, accountID, symbol string) (*types.AssetInfo, error)
	GetSchedule(ctx context.Context, symbol string) ([]types.Session, error)
	SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) bool) error
}
