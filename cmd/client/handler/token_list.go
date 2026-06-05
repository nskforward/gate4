package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/router"
)

func ListTokens(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requires 1 argument")
		}

		var userID string

		if len(args) == 1 {
			userID = args[0]
			args = args[1:]
		}

		tokens, err := c.TokenClient.ListUserTokens(ctx, userID)
		if err != nil {
			return err
		}

		if len(tokens) == 0 {
			fmt.Println("no tokens")
			return nil
		}
		fmt.Println()
		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println("| #  | ID                                   | CREATED    | EXPIRES    |")
		fmt.Println("-----------------------------------------------------------------------")

		for i, t := range tokens {
			fmt.Println(
				"|", fmt.Sprintf("%-2d", i+1),
				"|", t.ID,
				"|", t.Created.Format("2006-01-02"),
				"|", t.Expires.Format("2006-01-02"),
				"|")
		}

		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println()

		return nil
	}
}
