package cli

import (
	"context"
	"fmt"
)

func (c *Client) cmdPositions(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("requres 1 argument: account_key")
	}

	key := args[0]
	args = args[1:]

	positions, err := c.adminClient.Positions(ctx, key)
	if err != nil {
		return err
	}

	if len(positions) == 0 {
		fmt.Println("no positions")
		return nil
	}

	for i, pos := range positions {
		fmt.Printf("%d. %s at %s (%s)\n", i+1, pos.Symbol, pos.AveragePrice, pos.Size)
	}

	return nil
}
