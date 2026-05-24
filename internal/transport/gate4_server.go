package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/ssl"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Gate4Server struct {
	pb.UnimplementedGate4Server
	userStore     users.Store
	keychainStore *keychain.Store
	brokerPool    *brokers.Pool
}

func NewGate4Server(cfg config.Config) (*Gate4Server, error) {
	userStore, err := users.NewFileStorage()
	if err != nil {
		return nil, err
	}
	caKey, caCert, err := loadCA(cfg)
	if err != nil {
		return nil, err
	}
	return &Gate4Server{
		userStore:     userStore,
		keychainStore: keychain.NewStore(caKey, caCert),
		brokerPool:    brokers.NewPool(),
	}, nil
}

func (gate4 *Gate4Server) CreateCert(ctx context.Context, req *pb.CreateCertRequest) (*pb.CreateCertResponse, error) {
	if req.CommonName == "" {
		return nil, status.Error(codes.InvalidArgument, "common name cannot be empty")
	}
	key, err := ssl.ParsePrivateKey([]byte(req.PrivateKey))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	cert, err := gate4.keychainStore.Generate(req.CommonName, key)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	buf, err := ssl.MarshalCert(cert)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.CreateCertResponse{
		Cert: buf.String(),
	}, nil
}

func (gate4 *Gate4Server) FindUser(ctx context.Context, req *pb.UserID) (*pb.User, error) {
	user, err := gate4.userStore.Find(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return convertOutUser(user), nil
}

func (gate4 *Gate4Server) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.ListUsersResponse, error) {
	users, err := gate4.userStore.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ListUsersResponse{Users: convertOutUsers(users)}, nil
}

func (gate4 *Gate4Server) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user := convertInUser(req)
	err := gate4.userStore.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertOutUser(user), nil
}

func (gate4 *Gate4Server) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.EmptyMessage, error) {
	err := gate4.userStore.Delete(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EmptyMessage{}, nil
}

func (gate4 *Gate4Server) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.EmptyMessage, error) {
	err := gate4.userStore.Block(ctx, req.UserId, req.Blocked)
	return &pb.EmptyMessage{}, err
}

func (gate4 *Gate4Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := gate4.userStore.Update(ctx, req.UserId, req.Secret, time.Unix(req.Expires, 0))
	if err != nil {
		return nil, err
	}
	return convertOutUser(user), nil
}

func (gate4 *Gate4Server) SubscribeQuotes(req *pb.SymbolRequest, stream grpc.ServerStreamingServer[pb.Quote]) error {
	user, err := gate4.userStore.Find(stream.Context(), req.UserId)
	if err != nil {
		return fmt.Errorf("search user error: %w", err)
	}
	client, err := gate4.brokerPool.Get(user)
	if err != nil {
		return fmt.Errorf("search broker client error: %w", err)
	}
	return client.SubscribeQuotes(stream.Context(), func(q types.Quote) error {
		return stream.Send(&pb.Quote{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp,
			AskPrice:  q.AskPrice,
			BidPrice:  q.BidPrice,
			AskSize:   q.AskSize,
			BidSize:   q.BidSize,
		})
	})
}
