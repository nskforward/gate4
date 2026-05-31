package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console"
)

// user list [-blocked] [-active]
func ListUsers(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {

		_, activeArg := console.FindArg("-active", args)
		_, blockedArg := console.FindArg("-blocked", args)

		users, err := c.UserClient.ListUsers(ctx)
		if err != nil {
			return err
		}

		maxLenID := 0

		filtered := users[:0]
		for _, user := range users {
			if len(user.ID) > maxLenID {
				maxLenID = len(user.ID)
			}
			if user.Blocked && activeArg && !blockedArg {
				continue
			}
			if !user.Blocked && blockedArg && !activeArg {
				continue
			}
			filtered = append(filtered, user)
		}

		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Created.Before(filtered[j].Created)
		})

		if len(filtered) == 0 {
			fmt.Println("no users")
			return nil
		}

		maskID := fmt.Sprintf("%%-%ds", maxLenID)

		fmt.Println(strings.Repeat("-", 31+maxLenID))
		fmt.Println("| # |", fmt.Sprintf(maskID, "USER ID"), "| STATUS  | CREATED    |")
		fmt.Println(strings.Repeat("-", 31+maxLenID))

		for i, user := range filtered {
			fmt.Println(
				fmt.Sprintf("| %d |", i+1),
				console.FormatText(fmt.Sprintf(maskID, user.ID), console.White),
				"|",
				formatStatus(user.Blocked),
				"|", user.Created.Format("2006-01-02"), "|",
			)
		}
		fmt.Println(strings.Repeat("-", 31+maxLenID))
		return nil
	}
}

/*
func CreateUser(client *transport.GrpcClient) Handler {
	return func(ctx context.Context, args []string) error {
		var user model.User

		scanner := console.NewScanner()
		defer scanner.Close()

		input, err := scanner.Scan(ctx, "broker id", "", nil)
		if err != nil {
			return err
		}
		user.BrokerID = users.BrokerID(input)
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
*/
/*
func DeleteUser(client *transport.GrpcClient) Handler {
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

		fmt.Println(console.FormatText("WARNING!", console.Red, console.Bold), "user will be permanently removed")
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
*/
/*
func BlockUser(client *transport.GrpcClient) Handler {
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

		return errors.New("not implemented")

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
*/

/*
func UpdateUser(client *transport.GrpcClient) Handler {
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

		_ = secret
		_ = validUntil

		return errors.New("not implemented")

			err = client.UpdateUser(ctx, userID, secret, validUntil)
			if err != nil {
				return err
			}


		fmt.Println("success: user update with id", userID)

		return nil
	}
}
*/

func formatStatus(blocked bool) string {
	status := console.FormatText("active ", console.Green)
	if blocked {
		status = console.FormatText("blocked", console.Red)
	}
	return fmt.Sprintf("%-7s", status)
}
