package main

import (
	"context"
	"fmt"
)

func Help(context.Context, []string) error {
	padding := "   "
	mask := "%-51s"

	fmt.Println("---------------------------------------------------")
	fmt.Println()
	fmt.Println("GATE 4 CLI client v0.0.1")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "help"), "- (show all commands)")
	fmt.Println(padding, fmt.Sprintf(mask, "cert create"), "- (generate a client cert)")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "user list [-blocked] [-active]"), "- (show users)")
	fmt.Println(padding, fmt.Sprintf(mask, "user create"), "- (create a new user)")
	fmt.Println(padding, fmt.Sprintf(mask, "user delete <user_id>"), "- (delete a user by id)")
	fmt.Println(padding, fmt.Sprintf(mask, "user edit <user_id>"), "- (update the user detailes by id)")
	fmt.Println(padding, fmt.Sprintf(mask, "user block <user_id>"), "- (block/unblock a user by id)")
	fmt.Println()
	fmt.Println(padding, fmt.Sprintf(mask, "schedule [-symbol <symbol>] [-id <user_id>]"), "- (get symbol schedule sessions)")
	fmt.Println(padding, fmt.Sprintf(mask, "subscribe quotes [-symbol <symbol>] [-id <user_id>]"), "- (subscribe for realtime quotes)")
	fmt.Println(padding, fmt.Sprintf(mask, "subscribe trades [-symbol <symbol>] [-id <user_id>]"), "- (subscribe for realtime quotes)")
	fmt.Println(padding, fmt.Sprintf(mask, "subscribe positions <user_id>"), "- (subscribe for positions update)")
	fmt.Println()
	return nil
}
