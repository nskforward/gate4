package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/input"
)

func EditUser(c *client.Client) Handler {
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

		oldUser.Name, err = scanner.Scan(ctx, "name", oldUser.Name, nil)
		if err != nil {
			return err
		}

		oldUser.Email, err = scanner.Scan(ctx, "email", oldUser.Email, nil)
		if err != nil {
			return err
		}

		oldUser.Blocked, err = scanner.ScanBool(ctx, "blocked?", new(oldUser.Blocked), nil)
		if err != nil {
			return err
		}

		err = c.UserClient.UpdateUser(ctx, oldUser)
		if err != nil {
			return err
		}

		fmt.Println("success: user updated")

		return nil
	}
}
