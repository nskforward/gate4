package cli

import (
	"context"
	"fmt"
	"sort"
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

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Start < sessions[j].Start && sessions[i].End < sessions[j].End
	})

	for i, sess := range sessions {
		fmt.Printf("%d. %s at (%s - %s)\n", i+1, sess.Type, time.Unix(sess.Start, 0).Format(time.RFC3339), time.Unix(sess.End, 0).Format(time.RFC3339))
	}

	return nil
}
