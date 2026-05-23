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

		fmt.Println("Please enter the following fields:")

		scanner := console.NewScanner()
		defer scanner.Close()

		fmt.Print("- broker id: ")
		input, err := scanner.Scan(ctx)
		if err != nil {
			return err
		}
		user.BrokerID = brokers.BrokerID(input)
		err = user.BrokerID.Validate()
		if err != nil {
			return err
		}

		fmt.Print("- account id: ")
		user.AccountID, err = scanner.Scan(ctx)
		if err != nil {
			return err
		}

		user.Secret, err = scanner.ScanPassword(ctx, "secret")
		if err != nil {
			return err
		}

		fmt.Print("- valid until (YYYY-MM-DD HH:MM): ")
		user.ValidUntil, err = scanner.ScanTime(ctx, "2006-01-02 15:04")
		if err != nil {
			return err
		}

		fmt.Print("- blocked? (y/n): ")
		user.Blocked, err = scanner.ScanBool(ctx)
		if err != nil {
			return err
		}

		err = client.AddUser(ctx, &user)
		if err != nil {
			return err
		}

		fmt.Println("user created with id:", user.ID)

		return nil
	}
}

func DeleteUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {

		fmt.Println("Please enter the following fields:")

		scanner := console.NewScanner()
		defer scanner.Close()

		fmt.Print("- user id: ")
		userID, err := scanner.Scan(ctx)
		if err != nil {
			return err
		}

		err = client.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}
		fmt.Println("user deleted")
		return nil
	}
}

func BlockUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {

		fmt.Println("Please enter the following fields:")

		scanner := console.NewScanner()
		defer scanner.Close()

		fmt.Print("- user id: ")
		userID, err := scanner.Scan(ctx)
		if err != nil {
			return err
		}

		fmt.Print("- blocked? (y/n): ")
		blocked, err := scanner.ScanBool(ctx)
		if err != nil {
			return err
		}

		err = client.BlockUser(ctx, userID, blocked)
		if err != nil {
			return err
		}

		if blocked {
			fmt.Println("user blocked")
		} else {
			fmt.Println("user unblocked")
		}

		return nil
	}
}

func UpdateUser(client *grpc.AdminClient) Handler {
	return func(ctx context.Context, args []string) error {

		fmt.Println("Please enter the following fields:")

		scanner := console.NewScanner()
		defer scanner.Close()

		fmt.Print("- user id: ")
		userID, err := scanner.Scan(ctx)
		if err != nil {
			return err
		}

		secret, err := scanner.ScanPassword(ctx, "secret")
		if err != nil {
			return err
		}

		fmt.Print("- valid until (YYYY-MM-DD HH:MM): ")
		validUntil, err := scanner.ScanTime(ctx, "2006-01-02 15:04")
		if err != nil {
			return err
		}

		err = client.UpdateUser(ctx, userID, secret, validUntil)
		if err != nil {
			return err
		}

		fmt.Println("user update")

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
