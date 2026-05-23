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

func (c *AdminClient) CreateCert(ctx context.Context, commonName, privateKey string) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.CreateCert(reqCtx, &pb.CreateCertRequest{
		CommonName: commonName,
		PrivateKey: privateKey,
	})
	if err != nil {
		return "", err
	}
	return resp.Cert, nil
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

func (c *AdminClient) CreateUser(ctx context.Context, user *users.User) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.CreateUser(reqCtx, convertOutUser(user))
	if err != nil {
		return err
	}
	user.ID = *resp.Id
	user.Created = time.Unix(resp.Created, 0)
	return nil
}

func (c *AdminClient) DeleteUser(ctx context.Context, userID string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.DeleteUser(reqCtx, &pb.UserID{
		UserId: userID,
	})
	return err
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

func (c *AdminClient) UpdateUser(ctx context.Context, userID, secret string, expires time.Time) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.UpdateUser(reqCtx, &pb.UpdateUserRequest{
		UserId:  userID,
		Secret:  secret,
		Expires: expires.Unix(),
	})
	return err
}
