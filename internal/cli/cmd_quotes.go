package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) cmdQuotes(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("requres <symbol> <key> arguments")
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

		ask, err := decimal.NewFromString(resp.Ask)
		if err != nil {
			fmt.Println("error: bad ask format:", err)
			continue
		}

		bid, err := decimal.NewFromString(resp.Bid)
		if err != nil {
			fmt.Println("error: bad bid format:", err)
			continue
		}

		spread := ask.Sub(bid)

		fmt.Println(time.Unix(resp.Timestamp, 0).Format("2006-01-02 15:04"), resp.Symbol, "| ask", ask.StringFixed(2), "| bid", bid.StringFixed(2), "| spread", spread.StringFixed(2))
	}
}
