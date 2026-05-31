package server

import (
	"context"

	"github.com/nskforward/gate4/internal/api/grpc/common"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, err
	}
	return &pb.UserList{
		Users: common.ConvertInUsers(users),
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	modelUser := common.ConvertOutUser(user)
	err := h.userService.Create(ctx, &modelUser)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return common.ConvertInUser(modelUser), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.EmptyMessage, error) {
	return nil, h.userService.Delete(ctx, req.UserId)
}
