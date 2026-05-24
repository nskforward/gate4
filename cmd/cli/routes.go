package main

import (
	"github.com/nskforward/gate4/internal/cli"
	"github.com/nskforward/gate4/internal/transport"
)

func routes(r *cli.Router, grpcClient *transport.Gate4Client) {
	r.Handle("help", cli.Help)
	r.Handle("cert create", cli.CreateCert(grpcClient))
	r.Handle("user list", cli.ListUsers(grpcClient))
	r.Handle("user create", cli.CreateUser(grpcClient))
	r.Handle("user delete", cli.DeleteUser(grpcClient))
	r.Handle("user block", cli.BlockUser(grpcClient))
	r.Handle("user edit", cli.UpdateUser(grpcClient))
	r.Handle("subscribe quotes", cli.SubscribeQuotes(grpcClient))
}
