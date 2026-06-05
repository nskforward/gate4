package handler

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
	"github.com/nskforward/gate4/pkg/console/router"
)

func CreateUser(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {
		var user model.User
		var err error

		scanner := input.NewScanner()
		defer scanner.Close()

		user.Name, err = scanner.Scan(ctx, "name", "", nil)
		if err != nil {
			return err
		}

		fmt.Println()

		role, err := scanner.ScanOption(ctx, "choose the role", "", []input.Option{
			{Key: "screener", Description: "screener (read market data only)"},
			{Key: "trader", Description: "trader (read market data and place orders)"},
			{Key: "admin", Description: "admin (access management only)"},
		})
		if err != nil {
			return err
		}

		user.Role = model.Role(role)

		fmt.Println()

		user.Email, err = scanner.Scan(ctx, "email", "", nil)
		if err != nil {
			return err
		}

		fmt.Println()

		user.Blocked, err = scanner.ScanBool(ctx, "blocked?", new(false), nil)
		if err != nil {
			return err
		}

		fmt.Println()

		err = c.UserClient.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		fmt.Println(output.FormatText("user created:", output.Green))
		fmt.Println("- id:", user.ID)
		fmt.Println("- name:", user.Name)
		fmt.Println("- email:", user.Email)
		fmt.Println("- role:", user.Role.String())
		fmt.Println("- blocked:", user.Blocked)
		fmt.Println("- created:", user.Created.Format("2006-01-02 15:04:05"))
		fmt.Println()

		return nil
	}
}
