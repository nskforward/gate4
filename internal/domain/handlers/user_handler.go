package handlers

import (
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
}

func (h *UserHandler) Register(serv *grpc.Server) {
	pb.RegisterUserServiceServer(serv, h)
}
