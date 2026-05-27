package main

import (
	"github.com/nskforward/gate4/internal/app/client"
	"github.com/nskforward/gate4/internal/cli"
)

func routes(r *cli.Router, app *client.App) {
	r.Handle("help", cli.Help)
	r.Handle("cert create", cli.CreateCert(app))
	r.Handle("user list", cli.ListUsers(app))
	r.Handle("user create", cli.CreateUser(app))
	r.Handle("user delete", cli.DeleteUser(app))
	r.Handle("user block", cli.BlockUser(app))
	r.Handle("user edit", cli.UpdateUser(app))
	r.Handle("subscribe quotes", cli.SubscribeQuotes(app))
}
