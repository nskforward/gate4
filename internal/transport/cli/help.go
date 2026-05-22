package cli

import (
	"context"
	"fmt"
)

func Help(context.Context, []string) error {
	padding := "   "
	mask := "%-8s"

	fmt.Println("---------------------------------------------------")
	fmt.Println()
	fmt.Println("GATE 4 CLI client v0.0.1")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "help"), "- (show commands)")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "users"), "- (show all users)")
	fmt.Println(padding, fmt.Sprintf(mask, "user add"), "- (add user)")
	fmt.Println(padding, fmt.Sprintf(mask, "user del"), "- (delete user by id)")
	fmt.Println(padding, fmt.Sprintf(mask, "user get"), "- (search user by id)")
	fmt.Println()
	return nil
}
