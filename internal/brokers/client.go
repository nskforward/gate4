package brokers

import (
	"context"

	"github.com/nskforward/gate4/pkg/types"
)

type Client interface {
	SubscribeQuotes(context.Context, func(types.Quote) error) error
}
