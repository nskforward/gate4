package transport

import (
	"context"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gate4Client struct {
	client pb.Gate4Client
	conn   *grpc.ClientConn
}

func NewGate4Client(network, address string) (*Gate4Client, error) {
	normalizedAddr := address
	if network == "unix" {
		normalizedAddr = "unix:///" + filepath.ToSlash(address)
	}

	conn, err := grpc.NewClient(normalizedAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Gate4Client{
		conn:   conn,
		client: pb.NewGate4Client(conn),
	}, nil
}

func (c *Gate4Client) Close() {
	c.conn.Close()
}

func (c *Gate4Client) SubscribeQuotes(ctx context.Context, userID, symbol string, send func(quote types.Quote) error) error {
	stream, err := c.client.SubscribeQuotes(ctx, &pb.SymbolRequest{
		UserId: userID,
		Symbol: symbol,
	})
	if err != nil {
		return err
	}
	for {
		q, err := stream.Recv()
		if err != nil {
			return err
		}
		err = send(types.Quote{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp,
			Ask:       q.Ask,
			Bid:       q.Bid,
		})
		if err != nil {
			return err
		}
	}
}

func (c *Gate4Client) CreateCert(ctx context.Context, commonName, privateKey string) (string, error) {
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

func (c *Gate4Client) ListUsers(ctx context.Context) ([]*users.User, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.ListUsers(reqCtx, &pb.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	return convertInUsers(resp.Users), nil
}

func (c *Gate4Client) CreateUser(ctx context.Context, user *users.User) error {
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

func (c *Gate4Client) DeleteUser(ctx context.Context, userID string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.DeleteUser(reqCtx, &pb.UserID{
		UserId: userID,
	})
	return err
}

func (c *Gate4Client) BlockUser(ctx context.Context, userID string, blocked bool) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.BlockUser(reqCtx, &pb.BlockUserRequest{
		UserId:  userID,
		Blocked: blocked,
	})
	return err
}

func (c *Gate4Client) UpdateUser(ctx context.Context, userID, secret string, expires time.Time) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.UpdateUser(reqCtx, &pb.UpdateUserRequest{
		UserId:  userID,
		Secret:  secret,
		Expires: expires.Unix(),
	})
	return err
}
