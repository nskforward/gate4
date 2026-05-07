package types

import (
	"context"

	"github.com/nskforward/gate4/pkg/stream"
)

type BrokerClient interface {
	GetAccount(ctx context.Context, accountID string) (*AccountInfo, error)
	GetAsset(ctx context.Context, accountID, symbol string) (*AssetInfo, error)
	GetSchedule(ctx context.Context, symbol string) ([]Session, error)
	SubscribeQuotes(ctx context.Context, symbols []string) (*stream.Stream[Quote], error)
}
