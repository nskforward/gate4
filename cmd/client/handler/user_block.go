package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
)

func BlockUser(c *client.Client) Handler {
	return changeUserStatus(c, true)
}

func UnblockUser(c *client.Client) Handler {
	return changeUserStatus(c, false)
}

func changeUserStatus(c *client.Client, blockedAction bool) func(ctx context.Context, args []string) error {
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

		fmt.Println(
			output.FormatText("WARNING!", output.Yellow, output.Bold),
			fmt.Sprintf("the following user will be %s:", getStatusActionMessage(blockedAction)),
			fmt.Sprintf("%s (%s)", oldUser.Name, oldUser.Email),
		)

		allow, err := scanner.ScanBool(ctx, "continue?", nil, nil)
		if err != nil {
			return err
		}
		if !allow {
			fmt.Println("canceled")
			return nil
		}

		oldUser.Blocked = blockedAction

		err = c.UserClient.UpdateUser(ctx, oldUser)
		if err != nil {
			return err
		}

		if blockedAction {
			fmt.Println("success: user has been blocked")
		} else {
			fmt.Println("success: user has been unblocked")
		}

		return nil
	}
}

func getStatusActionMessage(defaultBlocked bool) string {
	if defaultBlocked {
		return output.FormatText("blocked", output.Red)
	}
	return output.FormatText("unblocked", output.Green)
}
