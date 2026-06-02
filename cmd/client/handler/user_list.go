package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/input"
	"github.com/nskforward/gate4/pkg/console/output"
)

// user list [-blocked] [-active]
func ListUsers(c *client.Client) Handler {
	return func(ctx context.Context, args []string) error {

		_, activeArg := input.FindArg("-active", args)
		_, blockedArg := input.FindArg("-blocked", args)

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

func formatStatus(blocked bool) string {
	status := output.FormatText("active ", output.Green)
	if blocked {
		status = output.FormatText("blocked", output.Red)
	}
	return fmt.Sprintf("%-7s", status)
}

func countDigits(n int) int {
	return len(fmt.Sprintf("%d", n))
}
