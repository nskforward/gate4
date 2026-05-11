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

	sessions, current, err := c.adminClient.GetSchedule(ctx, key, symbol)
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
		if sess.Start == current.Start && sess.End == current.End {
			fmt.Println("-----------------------------------------------------")
		}
		fmt.Printf("%02d. %s\t(%s - %s)\n", i+1, sess.Type, time.Unix(sess.Start, 0).Format("2006-01-02 15:04"), time.Unix(sess.End, 0).Format("2006-01-02 15:04"))
		if sess.Start == current.Start && sess.End == current.End {
			fmt.Println("-----------------------------------------------------")
		}
	}

	return nil
}
