package main

import (
	"github.com/nskforward/gate4/cmd/client/handler"
	"github.com/nskforward/gate4/internal/api/grpc/client"
)

func routes(r *handler.Router, c *client.Client) {
	r.Handle("help", Help)
	r.Handle("user list", handler.ListUsers(c))
	r.Handle("user create", handler.CreateUser(c))
	r.Handle("user delete", handler.DeleteUser(c))
	r.Handle("user edit", handler.UpdateUser(c))
	//r.Handle("cert create", handler.CreateCert(client))
	//r.Handle("user block", handler.BlockUser(client))
	//r.Handle("subscribe quotes", handler.SubscribeQuotes(client))
}
