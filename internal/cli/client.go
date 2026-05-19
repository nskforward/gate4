package cli

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/transport"
)

type Client struct {
	adminClient *transport.AdminClient
}

func NewClient(adminClient *transport.AdminClient) *Client {
	return &Client{
		adminClient: adminClient,
	}
}

func (c *Client) Run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("no command provided (run `gate4 help` to see possible commands)")
	}

	signalCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	path := args[0]
	_ = path
	command := args[1]
	args = args[2:]

	return c.execCmd(signalCtx, command, args)
}

func (c *Client) execCmd(ctx context.Context, command string, args []string) error {
	switch command {
	case "help":
		return c.cmdHelp(ctx)
	case "cert":
		return c.cmdCert(ctx, args)
	case "account":
		return c.cmdAccount(ctx, args)
	case "quotes":
		return c.cmdQuotes(ctx, args)
	case "position":
		return c.cmdPositions(ctx, args)
	case "schedule":
		return c.cmdSchedule(ctx, args)
	case "my-trades":
		return c.cmdMyTrades(ctx, args)
	case "asset":
		return c.cmdAssetInfo(ctx, args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}
