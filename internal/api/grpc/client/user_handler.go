package client

import (
	"context"
	"errors"
	"time"

	"github.com/nskforward/gate4/internal/api/grpc/common"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			return nil, errors.New("server unavailable")
		}
		return nil, err
	}
	return common.ConvertOutUsers(resp.Users), nil
}

func (h *UserHandler) CreateUser(ctx context.Context, user *model.User) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := h.client.CreateUser(reqCtx, common.ConvertInUser(*user))
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			return errors.New("server unavailable")
		}
		return err
	}
	*user = common.ConvertOutUser(resp)
	return nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, userID string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := h.client.DeleteUser(reqCtx, &pb.UserID{
		UserId: userID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			return errors.New("server unavailable")
		}
		return err
	}
	return nil
}
