package finam

import (
	"context"
	"time"

	"github.com/nskforward/gate4/pkg/types"
)

func (client *Client) SubscribeQuotes(ctx context.Context, symbol string, send func(types.Quote) error) error {
	token, err := client.GetToken()
	if err != nil {
		return err
	}

	ctx = token.Context(ctx)

	for range 10 {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
			send(types.Quote{Symbol: "TST", Timestamp: time.Now().Unix(), Ask: "0", Bid: "0"})
		}
	}
	return nil
}
