package transport

import (
	"os"
	"path/filepath"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	UserClient *client.UserHandler
	conn       *grpc.ClientConn
}

func NewGrpcClient() (*GrpcClient, error) {
	addr := "unix:///" + filepath.ToSlash(filepath.Join(os.TempDir(), "gate4.sock"))

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		UserClient: client.NewUserHandler(conn),
		conn:       conn,
	}, nil
}

func (c *GrpcClient) Close() {
	c.conn.Close()
}
