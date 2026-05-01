package cli

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if len(args) < 2 {
		return fmt.Errorf("requres <symbol> <account_key> arguments")
	}

	symbol := args[0]
	accountKey := args[1]
	args = args[2:]

	stream, err := c.adminClient.SubscribeQuotes(ctx, accountKey, symbol)
	if err != nil {
		return fmt.Errorf("gate4 server grpc error: %w", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("-- stream closed by server --")
			return nil
		}
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return err
			}
			code := st.Code()
			if code == codes.Canceled {
				return nil
			}
			return fmt.Errorf("stream closed: %s", st.Message())
		}
		fmt.Println("quote", resp.Timestamp, resp.BrokerId, resp.Symbol, resp.Ask, resp.Bid)
	}
}
