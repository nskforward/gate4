package types

import (
	"context"

	"github.com/nskforward/gate4/pkg/streams"
)

type BrokerClient interface {
	GetAccount(ctx context.Context, accountID string) (*AccountInfo, error)
	GetAsset(ctx context.Context, accountID, symbol string) (*AssetInfo, error)
	GetSchedule(ctx context.Context, symbol string) ([]Session, error)
	SubscribeQuotes(context.Context, string) *streams.Stream[Quote]
}
