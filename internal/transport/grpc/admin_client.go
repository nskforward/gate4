package grpc

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	transport "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AdminClient struct {
	client pb.AdminClient
	conn   *transport.ClientConn
}

func NewAdminClient(addr string) (*AdminClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &AdminClient{
		conn:   conn,
		client: pb.NewAdminClient(conn),
	}, nil
}

func (c *AdminClient) Close() {
	c.conn.Close()
}

func (c *AdminClient) ListUsers(ctx context.Context) ([]*users.User, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.ListUsers(reqCtx, &pb.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	return convertInUsers(resp.Users), nil
}

func (c *AdminClient) AddUser(ctx context.Context, user *users.User) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.AddUser(reqCtx, convertOutUser(user))
	if err != nil {
		return err
	}
	user.ID = resp.UserId
	return nil
}

func (c *AdminClient) BlockUser(ctx context.Context, userID string, blocked bool) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.BlockUser(reqCtx, &pb.BlockUserRequest{
		UserId:  userID,
		Blocked: blocked,
	})
	return err
}
