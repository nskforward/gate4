package client

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/domain/handler/grpc/converter"
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

func (h *UserHandler) ListUsers(ctx context.Context) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := h.client.ListUsers(ctx, &pb.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	return converter.ConvertOutUsers(resp.Users), nil
}

func (h *UserHandler) CreateUser(ctx context.Context, user *model.User) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := h.client.CreateUser(reqCtx, converter.ConvertInUser(*user))
	if err != nil {
		return err
	}
	*user = converter.ConvertOutUser(resp)
	return nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, userID string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := h.client.DeleteUser(reqCtx, &pb.UserID{
		UserId: userID,
	})
	return err
}
