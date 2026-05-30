package user

import (
	"context"

	"github.com/nskforward/gate4/internal/domain/model"
	usecases "github.com/nskforward/gate4/internal/domain/usecases/user"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	pb.UnimplementedUsersServer
	userUsecases *usecases.UserUsecases
}

func NewUserHandler(userUsecases *usecases.UserUsecases) *UserHandler {
	return &UserHandler{
		userUsecases: userUsecases,
	}
}

func (h *UserHandler) Register(serv *grpc.Server) {
	pb.RegisterUsersServer(serv, h)
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.UserList, error) {
	users, err := h.userUsecases.List(ctx)
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
