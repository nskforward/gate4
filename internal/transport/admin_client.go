package transport

import (
	"context"
	"time"

	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AdminClient struct {
	client pb.AdminClient
	conn   *grpc.ClientConn
}

func NewAdminClient(addr string) (*AdminClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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

func (c *AdminClient) ListAccounts(ctx context.Context) ([]*pb.Account, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.ListAccounts(reqCtx, &pb.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *AdminClient) AddAccount(ctx context.Context, brokerID, id, secret string, validDate int64) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.AddAccount(reqCtx, &pb.AddAccountRequest{
		Account: &pb.Account{
			Id:         id,
			BrokerId:   brokerID,
			Secret:     secret,
			ValidUntil: validDate,
		},
	})
	return err
}

func (c *AdminClient) DeleteAccount(ctx context.Context, key string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.DeleteAccount(reqCtx, &pb.DeleteAccountRequest{
		Key: key,
	})
	return err
}
