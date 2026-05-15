package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) cmdMyTrades(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requres <key> arguments")
	}

	accountKey := args[0]
	args = args[1:]

	stream, err := c.adminClient.SubscribeAccountTrades(ctx, accountKey)
	if err != nil {
		return fmt.Errorf("server error: %w", err)
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

		fmt.Println(time.Unix(resp.Timestamp, 0).Format("2006-01-02 15:04"), resp.Symbol, "| price", resp.Price, "| size", resp.Size)
	}
}
