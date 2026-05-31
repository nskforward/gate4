package server

import (
	"context"

	"github.com/nskforward/gate4/internal/api/grpc/converter"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	pb.UnimplementedUsersServer
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(servers ...*grpc.Server) {
	for _, s := range servers {
		pb.RegisterUsersServer(s, h)
	}
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.UserList, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		// TODO: handle not found
		return nil, err
	}
	return &pb.UserList{
		Users: converter.ConvertInUsers(users),
	}, nil
}
