package cli

import (
	"context"
	"fmt"
)

func Help(context.Context, []string) error {
	padding := "   "
	mask := "%-10s"

	fmt.Println("---------------------------------------------------")
	fmt.Println()
	fmt.Println("GATE 4 CLI client v0.0.1")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "help"), "- (show all commands)")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "show users"), "- (show all users)")
	fmt.Println(padding, fmt.Sprintf(mask, "create user"), "- (create new user)")
	fmt.Println(padding, fmt.Sprintf(mask, "delete user"), "- (delete user by id)")
	fmt.Println(padding, fmt.Sprintf(mask, "edit user"), "- (update user secret and valid date by id)")
	fmt.Println(padding, fmt.Sprintf(mask, "block user"), "- (block/unblock user by id)")
	fmt.Println()
	return nil
}
