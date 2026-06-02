package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/internal/domain/model"
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

		maxName := 0
		maxEmail := 0

		filtered := users[:0]
		for _, user := range users {
			if user.Blocked && activeArg && !blockedArg {
				continue
			}
			if !user.Blocked && blockedArg && !activeArg {
				continue
			}
			filtered = append(filtered, user)
			if len(user.Name) > maxName {
				maxName = len(user.Name)
			}
			if len(user.Email) > maxEmail {
				maxEmail = len(user.Email)
			}
		}

		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Created.Before(filtered[j].Created)
		})

		if len(filtered) == 0 {
			fmt.Println("no users")
			return nil
		}

		nameMask := fmt.Sprintf("%%-%ds", maxName)
		emailMask := fmt.Sprintf("%%-%ds", maxEmail)
		idMask := fmt.Sprintf("%%-%dv", countDigits(len(filtered)))

		fmt.Println(strings.Repeat("-", 73+maxName+maxEmail))
		fmt.Println(
			"|", fmt.Sprintf(idMask, "#"),
			"|", "USER ID                             ",
			"|", fmt.Sprintf(nameMask, "NAME"),
			"|", fmt.Sprintf(emailMask, "EMAIL"),
			"| STATUS  | CREATED    |")
		fmt.Println(strings.Repeat("-", 73+maxName+maxEmail))

		for i, user := range filtered {
			fmt.Println("|",
				fmt.Sprintf(idMask, i+1), "|",
				user.ID, "|",
				fmt.Sprintf(nameMask, user.Name), "|",
				fmt.Sprintf(emailMask, user.Email), "|",
				formatStatus(user.Blocked), "|",
				user.Created.Format("2006-01-02"), "|",
			)
		}
		fmt.Println(strings.Repeat("-", 73+maxName+maxEmail))
		return nil
	}
}

func CreateUser(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {
		var user model.User
		var err error

		scanner := console.NewScanner()
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

func DeleteUser(c *client.Client) Handler {
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

		fmt.Println(console.FormatText("WARNING!", console.Yellow, console.Bold), "user will be permanently removed")
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

func UpdateUser(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requeres 1 argument")
		}

		userID := args[0]
		args = args[1:]

		oldUser, err := c.UserClient.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		oldUser.Name, err = scanner.Scan(ctx, "name", oldUser.Name, nil)
		if err != nil {
			return err
		}

		oldUser.Email, err = scanner.Scan(ctx, "email", oldUser.Email, nil)
		if err != nil {
			return err
		}

		oldUser.Blocked, err = scanner.ScanBool(ctx, "blocked?", new(oldUser.Blocked), nil)
		if err != nil {
			return err
		}

		err = c.UserClient.UpdateUser(ctx, oldUser)
		if err != nil {
			return err
		}

		fmt.Println("success: user updated")

		return nil
	}
}

func BlockUser(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requeres 1 argument")
		}

		userID := args[0]
		args = args[1:]

		oldUser, err := c.UserClient.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		fmt.Println(console.FormatText("WARNING!", console.Yellow, console.Bold), fmt.Sprintf("the following user will be %s:", console.FormatText("blocked", console.Red)), fmt.Sprintf("%s (%s)", oldUser.Name, oldUser.Email))
		allow, err := scanner.ScanBool(ctx, "continue?", nil, nil)
		if err != nil {
			return err
		}
		if !allow {
			fmt.Println("canceled")
			return nil
		}

		oldUser.Blocked = true

		err = c.UserClient.UpdateUser(ctx, oldUser)
		if err != nil {
			return err
		}

		fmt.Println("success: user has been blocked")

		return nil
	}
}

func UnblockUser(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requeres 1 argument")
		}

		userID := args[0]
		args = args[1:]

		oldUser, err := c.UserClient.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		scanner := console.NewScanner()
		defer scanner.Close()

		fmt.Println(console.FormatText("WARNING!", console.Yellow, console.Bold), fmt.Sprintf("the following user will be %s:", console.FormatText("unblocked", console.Green)), fmt.Sprintf("%s (%s)", oldUser.Name, oldUser.Email))
		allow, err := scanner.ScanBool(ctx, "continue?", nil, nil)
		if err != nil {
			return err
		}
		if !allow {
			fmt.Println("canceled")
			return nil
		}

		oldUser.Blocked = false

		err = c.UserClient.UpdateUser(ctx, oldUser)
		if err != nil {
			return err
		}

		fmt.Println("success: user has been unblocked")

		return nil
	}
}

func formatStatus(blocked bool) string {
	status := console.FormatText("active ", console.Green)
	if blocked {
		status = console.FormatText("blocked", console.Red)
	}
	return fmt.Sprintf("%-7s", status)
}

func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	// Для отрицательных чисел можно взять модуль (если нужно)
	if n < 0 {
		n = -n
	}
	return int(math.Floor(math.Log10(float64(n)))) + 1
}
