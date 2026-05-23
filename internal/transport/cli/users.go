package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/transport/grpc"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/console"
)

func ListUsers(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, _ []string) error {
		users, err := client.ListUsers(ctx)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			fmt.Println("no users")
			return nil
		}

		maxLenID := 0
		maxLenBlocked := len("active")
		for _, user := range users {
			if len(user.ID) > maxLenID {
				maxLenID = len(user.ID)
			}
			if user.Blocked {
				maxLenBlocked = len("blocked")
			}
		}

		maskID := fmt.Sprintf("%%-%ds", maxLenID)
		maskBlocked := fmt.Sprintf("%%-%ds", maxLenBlocked)

		for i, user := range users {
			left := time.Until(user.ValidUntil)
			hours := left.Hours()
			days := 0
			if hours >= 24 {
				days = int(hours) / 24
			}
			fmt.Println(fmt.Sprintf("%d.", i+1), fmt.Sprintf(maskID, user.ID), "|", formatBlocked(maskBlocked, user.Blocked, user.ValidUntil), "|", fmt.Sprintf("days left: %d", days))
		}
		return nil
	}
}

func AddUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {
		var user users.User

		scanner := console.NewScanner()
		defer scanner.Close()

		input, err := scanner.Scan(ctx, "broker id", "", nil)
		if err != nil {
			return err
		}
		user.BrokerID = brokers.BrokerID(input)
		err = user.BrokerID.Validate()
		if err != nil {
			return err
		}

		user.AccountID, err = scanner.Scan(ctx, "account id", "", nil)
		if err != nil {
			return err
		}

		user.Secret, err = scanner.ScanPassword(ctx, "secret")
		if err != nil {
			return err
		}

		user.ValidUntil, err = scanner.ScanTime(ctx, "valid until", "2006-01-02 15:04", time.Now().AddDate(1, 0, 0), nil)
		if err != nil {
			return err
		}

		user.Blocked, err = scanner.ScanBool(ctx, "blocked?", false, nil)
		if err != nil {
			return err
		}

		err = client.AddUser(ctx, &user)
		if err != nil {
			return err
		}

		fmt.Println("success: user created with id:", user.ID)

		return nil
	}
}

func DeleteUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {

		var argUserID string

		if len(args) == 1 {
			argUserID = args[0]
			args = args[1:]
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		userID, err := scanner.Scan(ctx, "user id", "", &argUserID)
		if err != nil {
			return err
		}

		fmt.Println("WARNING! user will be permanently removed")
		allow, err := scanner.ScanBool(ctx, "continue?", false, nil)
		if err != nil {
			return err
		}

		if !allow {
			return nil
		}

		err = client.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}

		fmt.Println("success: user deleted with id:", userID)
		return nil
	}
}

func BlockUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {

		var argUserID string

		if len(args) == 1 {
			argUserID = args[0]
			args = args[1:]
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		userID, err := scanner.Scan(ctx, "user id", "", &argUserID)
		if err != nil {
			return err
		}

		blocked, err := scanner.ScanBool(ctx, "blocked?", true, nil)
		if err != nil {
			return err
		}

		err = client.BlockUser(ctx, userID, blocked)
		if err != nil {
			return err
		}

		op := "blocked"
		if !blocked {
			op = "unblocked"
		}

		fmt.Println("success: user", op, "with id:", userID)

		return nil
	}
}

func UpdateUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {

		var argUserID string

		if len(args) == 1 {
			argUserID = args[0]
			args = args[1:]
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		userID, err := scanner.Scan(ctx, "user id", "", &argUserID)
		if err != nil {
			return err
		}

		secret, err := scanner.ScanPassword(ctx, "secret")
		if err != nil {
			return err
		}

		validUntil, err := scanner.ScanTime(ctx, "valid until", "2006-01-02 15:04", time.Now().AddDate(1, 0, 0), nil)
		if err != nil {
			return err
		}

		err = client.UpdateUser(ctx, userID, secret, validUntil)
		if err != nil {
			return err
		}

		fmt.Println("success: user update with id:", userID)

		return nil
	}
}

func formatBlocked(mask string, blocked bool, expires time.Time) string {
	if blocked {
		return fmt.Sprintf(mask, "blocked")
	}
	if time.Since(expires) > 0 {
		return fmt.Sprintf(mask, "blocked")
	}
	return fmt.Sprintf(mask, "active")
}
