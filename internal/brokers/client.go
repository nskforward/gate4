package brokers

import (
	"context"

	"github.com/nskforward/gate4/pkg/types"
)

type Client interface {
	Close() error
	SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error
}
