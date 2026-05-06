package cli

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) cmdCurrentSession(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("requres 2 argument: symbol, account_key")
	}

	symbol := args[0]
	key := args[1]

	args = args[2:]

	sess, err := c.adminClient.GetCurrentSession(ctx, key, symbol)
	if err != nil {
		return err
	}

	fmt.Printf("%s at (%s - %s)\n", sess.Type, time.Unix(sess.Start, 0).Format(time.RFC3339), time.Unix(sess.End, 0).Format(time.RFC3339))

	return nil
}
