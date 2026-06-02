package handler

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
	"github.com/nskforward/gate4/pkg/console/router"
)

func DeleteUser(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {

		var argUserID string

		if len(args) == 1 {
			argUserID = args[0]
			args = args[1:]
		}

		scanner := input.NewScanner()
		defer scanner.Close()

		userID, err := scanner.Scan(ctx, "user id", "", &argUserID)
		if err != nil {
			return err
		}

		fmt.Println(output.FormatText("WARNING!", output.Yellow, output.Bold), "user will be permanently removed")
		allow, err := scanner.ScanBool(ctx, "continue?", nil, nil)
		if err != nil {
			return err
		}

		if !allow {
			fmt.Println("canceled")
			return nil
		}

		err = c.UserClient.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}

		fmt.Println("success: user deleted", userID)
		return nil
	}
}
