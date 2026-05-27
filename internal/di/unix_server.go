package di

import (
	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

func NewUnixServer(apiServer *api.Server) *grpc.Server {
	unixSocket := grpc.NewServer()
	pb.RegisterGate4Server(unixSocket, apiServer)
	return unixSocket
}
