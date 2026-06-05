package client

import (
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	UserClient  *UserHandler
	TokenClient *TokenHandler
	conn        *grpc.ClientConn
}

func NewClient() (*Client, error) {
	addr := "unix:///" + filepath.ToSlash(filepath.Join(os.TempDir(), "gate4.sock"))

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		UserClient:  NewUserHandler(conn),
		TokenClient: NewTokenHandler(conn),
		conn:        conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
