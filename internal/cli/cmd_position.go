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

	list, err := c.adminClient.GetPositions(ctx, key)
	if err != nil {
		return err
	}

	if len(list.Positions) == 0 {
		fmt.Println("no positions")
		return nil
	}

	for i, pos := range list.Positions {
		size, err := decimal.NewFromString(pos.Size)
		if err != nil {
			fmt.Println("error: bad format of size:", size)
			return nil
		}

		direction := "long"
		if size.IsNegative() {
			direction = "short"
		}

		fmt.Println(fmt.Sprintf("%02d.", i+1), pos.Symbol, direction, "| price", pos.Price, "| size", size.StringFixed(2), "| profit", pos.Profit)
	}

	return nil
}

func (c *Client) cmdPosition(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("requres 2 argument: symbol, account_key")
	}

	symbol := args[0]
	key := args[1]
	args = args[2:]

	pos, err := c.adminClient.GetPosition(ctx, key, symbol)
	if err != nil {
		return err
	}

	if pos == nil {
		fmt.Println("no position")
		return nil
	}

	size, err := decimal.NewFromString(pos.Size)
	if err != nil {
		fmt.Println("error: bad format of size:", size)
		return nil
	}

	direction := "long"
	if size.IsNegative() {
		direction = "short"
	}

	fmt.Println(pos.Symbol, direction, "| price", pos.Price, "| size", size.StringFixed(2), "| profit", pos.Profit)

	return nil
}
