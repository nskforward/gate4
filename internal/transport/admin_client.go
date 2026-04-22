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

func (c *AdminClient) ListBrokers(ctx context.Context) ([]*pb.Broker, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.ListBrokers(reqCtx, &pb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
