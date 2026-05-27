package main

import (
	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/cli"
)

func routes(r *cli.Router, client *api.Client) {
	r.Handle("help", cli.Help)
	r.Handle("cert create", cli.CreateCert(client))
	r.Handle("user list", cli.ListUsers(client))
	r.Handle("user create", cli.CreateUser(client))
	r.Handle("user delete", cli.DeleteUser(client))
	r.Handle("user block", cli.BlockUser(client))
	r.Handle("user edit", cli.UpdateUser(client))
	r.Handle("subscribe quotes", cli.SubscribeQuotes(client))
}
