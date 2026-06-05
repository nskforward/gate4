package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
	"github.com/nskforward/gate4/pkg/console/router"
)

func CreateToken(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {
		if len(args) < 1 {
			return errors.New("requires 1 argument")
		}

		var userID string

		if len(args) == 1 {
			userID = args[0]
			args = args[1:]
		}

		token := model.Token{
			UserID: userID,
		}
		var err error

		scanner := input.NewScanner()
		defer scanner.Close()

		token.Expires, err = scanner.ScanTime(ctx, "expires", "2006-01-02", time.Now().AddDate(1, 0, 0), nil)
		if err != nil {
			return err
		}

		fmt.Println()

		err = c.TokenClient.CreateToken(ctx, &token)
		if err != nil {
			return err
		}

		fmt.Println(output.FormatText("token created:", output.Green))
		fmt.Println("- id:", token.ID)
		fmt.Println("- user_id:", token.UserID)
		fmt.Println("- created:", token.Created.Format("2006-01-02 15:04:05"))
		fmt.Println("- expires:", token.Expires.Format("2006-01-02 15:04:05"))
		fmt.Println()

		return nil
	}
}
