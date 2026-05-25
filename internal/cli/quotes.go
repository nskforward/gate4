package cli

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/types"
)

// subscribe quotes [-symbol <symbol>] [-id <id>]
func SubscribeQuotes(client *transport.Gate4Client) Handler {
	return func(ctx context.Context, args []string) error {

		argSymbol, _ := console.FindArg("-symbol", args)
		argID, _ := console.FindArg("-id", args)

		scanner := console.NewScanner()
		defer scanner.Close()

		symbol, err := scanner.Scan(ctx, "symbol", "", &argSymbol)
		if err != nil {
			return err
		}

		userID, err := scanner.Scan(ctx, "user id", "", &argID)
		if err != nil {
			return err
		}

		return client.SubscribeQuotes(ctx, userID, symbol, func(quote types.Quote) error {
			fmt.Println(quote.Timestamp, quote.Symbol, quote.Ask, quote.Bid)
			return nil
		})
	}
}
