package brokers

import (
	"context"

	"github.com/nskforward/gate4/pkg/types"
)

type Client interface {
	SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error
}
