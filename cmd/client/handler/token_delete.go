package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
	"github.com/nskforward/gate4/pkg/console/router"
)

func DeleteToken(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requires 1 argument")
		}

		var tokenID string

		if len(args) == 1 {
			tokenID = args[0]
			args = args[1:]
		}

		scanner := input.NewScanner()
		defer scanner.Close()

		fmt.Println(output.FormatText("WARNING!", output.Yellow, output.Bold), "token will be permanently deleted")
		allow, err := scanner.ScanBool(ctx, "continue?", nil, nil)
		if err != nil {
			return err
		}

		if !allow {
			fmt.Println("canceled")
			return nil
		}

		err = c.TokenClient.DeleteToken(ctx, tokenID)
		if err != nil {
			return err
		}

		fmt.Println("token deleted", tokenID)
		return nil
	}
}
