package cli

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

func (c *Client) cmdPositions(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("requres 1 argument: account_key")
	}

	key := args[0]
	args = args[1:]

	positions, err := c.adminClient.GetPositions(ctx, key)
	if err != nil {
		return err
	}

	if len(positions) == 0 {
		fmt.Println("no positions")
		return nil
	}

	posMask := "%d."
	if len(positions) > 9 {
		posMask = "%02d"
	}

	for i, pos := range positions {
		size, err := decimal.NewFromString(pos.Size)
		if err != nil {
			fmt.Println("error: bad format of size:", size)
			continue
		}
		direction := "long"
		if size.IsNegative() {
			direction = "short"
		}
		fmt.Println(fmt.Sprintf(posMask, i+1), pos.Symbol, direction, "| price", pos.Price, "| size", size.StringFixed(2), "| profit", pos.Profit)
	}

	return nil
}
