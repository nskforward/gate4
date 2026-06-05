package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
	"github.com/nskforward/gate4/pkg/console/router"
)

func EditUser(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requires 1 argument")
		}

		userID := args[0]
		args = args[1:]

		oldUser, err := c.UserClient.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		scanner := input.NewScanner()
		defer scanner.Close()

		fmt.Println()

		oldUser.Name, err = scanner.Scan(ctx, "name", oldUser.Name, nil)
		if err != nil {
			return err
		}

		fmt.Println()

		oldUser.Email, err = scanner.Scan(ctx, "email", oldUser.Email, nil)
		if err != nil {
			return err
		}

		fmt.Println()

		role, err := scanner.ScanOption(ctx, "role", oldUser.Role.String(), []input.Option{
			{Key: "screener", Description: "screener (read market data only)"},
			{Key: "trader", Description: "trader (read market data and place orders)"},
			{Key: "admin", Description: "admin (access management only)"},
		})
		if err != nil {
			return err
		}
		oldUser.Role = model.Role(role)

		fmt.Println()

		oldUser.Blocked, err = scanner.ScanBool(ctx, "blocked?", new(oldUser.Blocked), nil)
		if err != nil {
			return err
		}

		fmt.Println()

		err = c.UserClient.UpdateUser(ctx, oldUser)
		if err != nil {
			return err
		}

		fmt.Println(output.FormatText("user updated", output.Green))

		fmt.Println()

		return nil
	}
}
