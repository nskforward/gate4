package main

import (
	"github.com/nskforward/gate4/cmd/client/handler"
	"github.com/nskforward/gate4/internal/transport"
)

func routes(r *handler.Router, client *transport.GrpcClient) {
	r.Handle("help", Help)
	//r.Handle("cert create", handler.CreateCert(client))
	r.Handle("user list", handler.ListUsers(client))
	//r.Handle("user create", handler.CreateUser(client))
	//r.Handle("user delete", handler.DeleteUser(client))
	//r.Handle("user block", handler.BlockUser(client))
	//r.Handle("user edit", handler.UpdateUser(client))
	//r.Handle("subscribe quotes", handler.SubscribeQuotes(client))
}
