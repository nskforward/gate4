package cli

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/types"
)

func SubscribeQuotes(client *transport.Gate4Client) Handler {
	return func(ctx context.Context, args []string) error {

		var argSymbol, argUserID string

		if len(args) == 1 {
			argSymbol = args[0]
			args = args[1:]
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		symbol, err := scanner.Scan(ctx, "symbol", "", &argSymbol)
		if err != nil {
			return err
		}

		userID, err := scanner.Scan(ctx, "user id", "", &argUserID)
		if err != nil {
			return err
		}

		return client.SubscribeQuotes(ctx, userID, symbol, func(quote types.Quote) error {
			ask := "-"
			if len(quote.AskPrice) > 0 {
				ask = quote.AskPrice[0]
			}
			bid := "-"
			if len(quote.BidPrice) > 0 {
				bid = quote.BidPrice[0]
			}

			fmt.Println(quote.Timestamp, quote.Symbol, ask, bid)

			return nil
		})
	}
}
