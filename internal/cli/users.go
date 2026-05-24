package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/console"
)

func ListUsers(client *transport.Gate4Client) Handler {
	return func(ctx context.Context, args []string) error {

		filterActive := true
		filterBlocked := true

		for _, arg := range args {
			if arg == "-blocked" {
				filterBlocked = true
				filterActive = false
			}
			if arg == "-active" {
				filterBlocked = false
				filterActive = true
			}
		}

		users, err := client.ListUsers(ctx)
		if err != nil {
			return err
		}

		maxLenID := 0

		filtered := users[:0]
		for _, user := range users {
			if len(user.ID) > maxLenID {
				maxLenID = len(user.ID)
			}
			if user.Blocked && !filterBlocked {
				continue
			}
			if !user.Blocked && !filterActive {
				continue
			}
			filtered = append(filtered, user)
		}

		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Expires.Before(filtered[j].Expires)
		})

		if len(filtered) == 0 {
			fmt.Println("no users")
			return nil
		}

		maskID := fmt.Sprintf("%%-%ds", maxLenID)

		fmt.Println(strings.Repeat("-", 56+maxLenID))
		fmt.Println("| # |", fmt.Sprintf(maskID, "USER ID"), "| STATUS  | CREATED    | EXPIRES    | DAYS LEFT |")
		fmt.Println(strings.Repeat("-", 56+maxLenID))

		for i, user := range filtered {
			fmt.Println(
				fmt.Sprintf("| %d |", i+1),
				fmt.Sprintf(maskID, user.ID),
				"|",
				formatStatus(user.Blocked, user.Expires),
				"|", user.Created.Format("2006-01-02"),
				"|", user.Expires.Format("2006-01-02"),
				"|", formatDaysLeft(user.Expires, 14), "|",
			)
		}
		fmt.Println(strings.Repeat("-", 56+maxLenID))
		return nil
	}
}

func CreateUser(client *transport.Gate4Client) Handler {
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

		user.Expires, err = scanner.ScanTime(ctx, "expires", "2006-01-02 15:04", time.Now().AddDate(1, 0, 0), nil)
		if err != nil {
			return err
		}

		user.Blocked, err = scanner.ScanBool(ctx, "blocked?", new(false), nil)
		if err != nil {
			return err
		}

		err = client.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		fmt.Println("success: user created with id", user.ID)

		return nil
	}
}

func DeleteUser(client *transport.Gate4Client) Handler {
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
		allow, err := scanner.ScanBool(ctx, "continue?", nil, nil)
		if err != nil {
			return err
		}

		if !allow {
			fmt.Println("canceled")
			return nil
		}

		err = client.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}

		fmt.Println("success: user deleted with id", userID)
		return nil
	}
}

func BlockUser(client *transport.Gate4Client) Handler {
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

		blocked, err := scanner.ScanBool(ctx, "blocked?", nil, nil)
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

		fmt.Println("success: user", op, "with id", userID)

		return nil
	}
}

func UpdateUser(client *transport.Gate4Client) Handler {
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

		fmt.Println("success: user update with id", userID)

		return nil
	}
}

func formatStatus(blocked bool, expires time.Time) string {
	status := console.FormatText("active ", console.Green)
	if blocked || time.Since(expires) > 0 {
		status = console.FormatText("blocked", console.Red)
	}
	return fmt.Sprintf("%-7s", status)
}

func formatDaysLeft(expires time.Time, threshold int) string {
	daysLeft := int(time.Until(expires).Hours() / 24)
	text := fmt.Sprintf("%-9d", daysLeft)
	if daysLeft < threshold {
		return console.FormatText(text, console.Red)
	}
	return text
}
