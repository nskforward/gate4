package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/router"
)

func WhoAmI(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {
		if len(args) < 1 {
			return errors.New("requires 1 argument")
		}

		var tokenID string

		if len(args) == 1 {
			tokenID = args[0]
			args = args[1:]
		}

		user, err := c.TokenClient.Whoami(ctx, tokenID)
		if err != nil {
			return err
		}

		fmt.Println("---------------- USER -----------------")
		fmt.Println("name    :", user.Name)
		fmt.Println("email   :", user.Email)
		fmt.Println("role    :", user.Role.String())
		fmt.Println("blocked :", user.Blocked)
		fmt.Println("created :", user.Created)
		fmt.Println("---------------------------------------")

		return nil
	}
}
