package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/nskforward/gate4/internal/transport/grpc"
)

func ListUsers(client *grpc.AdminClient) Handler {

	formatBlocked := func(blocked bool, expires time.Time) string {
		if blocked {
			return "blocked"
		}
		if time.Since(expires) > 0 {
			return "blocked"
		}
		return "active"
	}

	return func(ctx context.Context, _ []string) error {
		users, err := client.ListUsers(ctx)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			fmt.Println("no users")
			return nil
		}
		for i, user := range users {
			left := time.Until(user.ValidUntil)
			hours := left.Hours()
			days := 0
			if hours >= 24 {
				days = int(hours) / 24
			}
			fmt.Println(fmt.Sprintf("%d.", i+1), user.ID, "|", formatBlocked(user.Blocked, user.ValidUntil), "|", fmt.Sprintf("days left: %d", days))
		}
		return nil
	}
}

func AddUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {
		return fmt.Errorf("add user: not implemented")
	}
}

func DeleteUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {
		return fmt.Errorf("not implemented")
	}
}
