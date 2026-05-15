package cli

import (
	"context"
	"fmt"
)

func (c *Client) cmdAssetInfo(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("requres 2 arguments: <symbol>, <account_key>")
	}

	symbol := args[0]
	key := args[1]
	args = args[2:]

	info, err := c.adminClient.GetAsset(ctx, key, symbol)
	if err != nil {
		return err
	}
	fmt.Println("symbol:", info.Symbol)
	fmt.Println("description:", info.Description)
	fmt.Println("currency:", info.Currency)
	fmt.Println("lot size:", info.LotSize)
	fmt.Println("min step:", info.MinStep)
	fmt.Println("decimals:", info.Decimals)

	return nil
}
