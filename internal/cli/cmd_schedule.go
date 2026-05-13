package cli

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) cmdSchedule(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("requres 2 argument: symbol, account_key")
	}

	symbol := args[0]
	key := args[1]

	args = args[2:]

	sessions, err := c.adminClient.GetSchedule(ctx, key, symbol)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("no sessions")
		return nil
	}

	for i, sess := range sessions {
		fmt.Printf("%02d. %s\t(%s - %s)\n", i+1, sess.Type, time.Unix(sess.Start, 0).Format("2006-01-02 15:04"), time.Unix(sess.End, 0).Format("2006-01-02 15:04"))
	}

	return nil
}
