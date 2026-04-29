package cli

import (
	"context"
	"fmt"
)

func (c *Client) cmdSubscribe(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("requres at least 1 argument")
	}

	command := args[0]
	args = args[1:]

	switch command {
	case "quotes":
		return c.cmdSubscribeQuotes(ctx, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (c *Client) cmdSubscribeQuotes(ctx context.Context, args []string) error {
	return fmt.Errorf("not implemented")
}
