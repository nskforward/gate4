package handler

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/internal/domain/usecase"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	pb.UnimplementedUsersServer
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(serv *grpc.Server) {
	pb.RegisterUsersServer(serv, h)
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.UserList, error) {
	users, err := h.userUsecase.List(ctx)
	// TODO: handle not found
	if err != nil {
		return nil, err
	}
	return &pb.UserList{
		Users: convertUsers(users),
	}, nil
}

func convertUsers(users []model.User) []*pb.User {
	result := make([]*pb.User, 0, len(users))
	for _, user := range users {
		result = append(result, &pb.User{
			Id: user.ID,
		})
	}
	return result
}
