package main

import (
	"context"
	"os"

	"github.com/nskforward/gate4/pkg/help"
)

func Help(context.Context, []string) error {
	menu := help.NewMenu("GATE 4 CLI client v0.0.1")

	userSection := menu.AddSection("user")
	userSection.AddCommand("user list [-blocked] [-active]", "show all users")
	userSection.AddCommand("user create", "create a new user")
	userSection.AddCommand("user delete <user_id>", "delete a user by id")
	userSection.AddCommand("user edit <user_id>", "update the user details by id")
	userSection.AddCommand("user block <user_id>", "block a user by id")
	userSection.AddCommand("user unblock <user_id>", "unblock a user by id")
	userSection.AddCommand("user cert <user_id>", "generate a cert for grpc client")

	accountSection := menu.AddSection("account")
	accountSection.AddCommand("account list <user_id>", "show the user accounts")
	accountSection.AddCommand("account add <user_id>", "add a new user account")
	accountSection.AddCommand("account delete <account_id>", "delete account by id")
	accountSection.AddCommand("account edit <account_id>", "update account by id")
	accountSection.AddCommand("account block <account_id>", "block an account by id")
	accountSection.AddCommand("account unblock <account_id>", "unblock an account by id")

	tradingSection := menu.AddSection("trading")
	tradingSection.AddCommand("schedule [-account_id] [-symbol]", "get the schedule")
	tradingSection.AddCommand("subscribe quotes [-account_id] [-symbol]", "subscribe for quotes")
	tradingSection.AddCommand("subscribe trades [-account_id] [-symbol]", "subscribe for trades")
	tradingSection.AddCommand("subscribe positions <account_id>", "subscribe for positions")

	_, err := menu.WriteTo(os.Stdout)
	return err
}
