package client

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/app/server"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type App struct {
	client pb.Gate4Client
	conn   *grpc.ClientConn
}

func NewApp() (*App, error) {
	addr := "unix:///" + filepath.ToSlash(server.UnixSocketPath)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &App{
		conn:   conn,
		client: pb.NewGate4Client(conn),
	}, nil
}

func (c *App) Close() {
	c.conn.Close()
}

func (c *App) SubscribeQuotes(ctx context.Context, userID, symbol string, send func(quote types.Quote) error) error {
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
			if errors.Is(err, io.EOF) {
				return nil
			}
			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.Canceled {
					return nil
				}
			}
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

func (c *App) CreateCert(ctx context.Context, commonName, privateKey string) (string, error) {
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

func (c *App) ListUsers(ctx context.Context) ([]*users.User, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.ListUsers(reqCtx, &pb.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	return api.ConvertInUsers(resp.Users), nil
}

func (c *App) CreateUser(ctx context.Context, user *users.User) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.CreateUser(reqCtx, api.ConvertOutUser(user))
	if err != nil {
		return err
	}
	user.ID = *resp.Id
	user.Created = time.Unix(resp.Created, 0)
	return nil
}

func (c *App) DeleteUser(ctx context.Context, userID string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.DeleteUser(reqCtx, &pb.UserID{
		UserId: userID,
	})
	return err
}

func (c *App) BlockUser(ctx context.Context, userID string, blocked bool) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.BlockUser(reqCtx, &pb.BlockUserRequest{
		UserId:  userID,
		Blocked: blocked,
	})
	return err
}

func (c *App) UpdateUser(ctx context.Context, userID, secret string, expires time.Time) error {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := c.client.UpdateUser(reqCtx, &pb.UpdateUserRequest{
		UserId:  userID,
		Secret:  secret,
		Expires: expires.Unix(),
	})
	return err
}
