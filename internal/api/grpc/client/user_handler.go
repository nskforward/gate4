package client

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/api/grpc/common"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	client pb.UsersClient
}

func NewUserHandler(conn *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		client: pb.NewUsersClient(conn),
	}
}

func (h *UserHandler) FindByID(ctx context.Context, userID string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	user, err := h.client.FindUserByID(ctx, &pb.UserID{UserId: userID})
	if err != nil {
		return model.User{}, wrapError(err)
	}
	return common.ConvertOutUser(user), nil
}

func (h *UserHandler) ListUsers(ctx context.Context) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := h.client.ListUsers(ctx, &pb.EmptyMessage{})
	if err != nil {
		return nil, wrapError(err)
	}
	return common.ConvertOutUsers(resp.Users), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := h.client.UpdateUser(ctx, common.ConvertInUser(user))
	return wrapError(err)
}

func (h *UserHandler) CreateUser(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := h.client.CreateUser(ctx, common.ConvertInUser(*user))
	if err != nil {
		return wrapError(err)
	}
	*user = common.ConvertOutUser(resp)
	return nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := h.client.DeleteUser(ctx, &pb.UserID{
		UserId: userID,
	})
	return wrapError(err)
}
