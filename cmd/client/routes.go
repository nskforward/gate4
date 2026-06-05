package main

import (
	"github.com/nskforward/gate4/cmd/client/handler"
	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/router"
)

func routes(r *router.Router, c *client.Client) {

	r.Handle("help", "show help menu", r.PrintHelp())
	r.Handle("whoami <token_id>", "check token", handler.WhoAmI(c))

	r.Handle("user list [-a, -b]", "show users with filter: -a active, -b blocked", handler.ListUsers(c))
	r.Handle("user create", "create a new user", handler.CreateUser(c))
	r.Handle("user delete <id>", "delete user", handler.DeleteUser(c))
	r.Handle("user edit <id>", "change user details", handler.EditUser(c))
	r.Handle("user block <id>", "block a user", handler.BlockUser(c))
	r.Handle("user unblock <id>", "unblock a user", handler.UnblockUser(c))

	r.Handle("token list <user_id>", "show user tokens", handler.ListTokens(c))
	r.Handle("token create <user_id>", "create user token", handler.CreateToken(c))
	r.Handle("token delete <token_id>", "delete token", handler.DeleteToken(c))

	//r.Handle("cert create", handler.CreateCert(client))
	//r.Handle("subscribe quotes", handler.SubscribeQuotes(client))
}

/*
	tradingSection := menu.AddSection("trading")
	tradingSection.AddCommand("schedule [-account_id] [-symbol]", "get the schedule")
	tradingSection.AddCommand("subscribe quotes [-account_id] [-symbol]", "subscribe for quotes")
	tradingSection.AddCommand("subscribe trades [-account_id] [-symbol]", "subscribe for trades")
	tradingSection.AddCommand("subscribe positions <account_id>", "subscribe for positions")
*/
