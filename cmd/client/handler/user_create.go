package handler

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/console/input"
)

func CreateUser(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {
		var user model.User
		var err error

		scanner := input.NewScanner()
		defer scanner.Close()

		user.Name, err = scanner.Scan(ctx, "name", "", nil)
		if err != nil {
			return err
		}

		user.Email, err = scanner.Scan(ctx, "email", "", nil)
		if err != nil {
			return err
		}

		user.Blocked, err = scanner.ScanBool(ctx, "blocked?", new(false), nil)
		if err != nil {
			return err
		}

		err = c.UserClient.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		fmt.Println("success: user created")
		fmt.Println()
		fmt.Println("user detailes:")
		fmt.Println("- id:", user.ID)
		fmt.Println("- name:", user.Name)
		fmt.Println("- email:", user.Email)
		fmt.Println("- blocked:", user.Blocked)
		fmt.Println("- created:", user.Created.Format("2006-01-02 15:04:05"))

		return nil
	}
}
