package handler

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type GRPCUserHandler struct {
	pb.UnimplementedUsersServer
	userService *service.UserService
}

func NewGRPCUserHandler(userService *service.UserService) *GRPCUserHandler {
	return &GRPCUserHandler{
		userService: userService,
	}
}

func (h *GRPCUserHandler) Register(servers ...*grpc.Server) {
	for _, s := range servers {
		pb.RegisterUsersServer(s, h)
	}
}

func (h *GRPCUserHandler) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.UserList, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		// TODO: handle not found
		return nil, err
	}
	return &pb.UserList{
		Users: convertOutUsers(users),
	}, nil
}

func convertOutUsers(users []model.User) []*pb.User {
	result := make([]*pb.User, 0, len(users))
	for _, user := range users {
		result = append(result, &pb.User{
			Id: user.ID,
		})
	}
	return result
}
