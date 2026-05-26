package transport

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/ssl"
	"github.com/nskforward/gate4/pkg/tools"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Gate4Server struct {
	pb.UnimplementedGate4Server
	userStore     users.Store
	keychainStore *keychain.Store
	clientPool    *brokers.ClientPool
	serverCtx     context.Context
	cancel        context.CancelFunc
}

func NewGate4Server(ctx context.Context, cfg config.Config) (*Gate4Server, error) {
	userStore, err := users.NewFileStorage()
	if err != nil {
		return nil, err
	}

	caKey, caCert, err := loadCA(cfg)
	if err != nil {
		return nil, err
	}

	serverCtx, cancel := context.WithCancel(ctx)

	return &Gate4Server{
		serverCtx:     serverCtx,
		cancel:        cancel,
		userStore:     userStore,
		keychainStore: keychain.NewStore(caKey, caCert),
		clientPool:    brokers.NewClientPool(),
	}, nil
}

func (gate4 *Gate4Server) Close() {
	gate4.cancel()
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

	_, err := gate4.clientPool.GetOrCreate(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = gate4.userStore.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertOutUser(user), nil
}

func (gate4 *Gate4Server) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.EmptyMessage, error) {
	user, err := gate4.userStore.Find(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = gate4.clientPool.Delete(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = gate4.userStore.Delete(ctx, req.UserId)
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

	ctx, cancel := tools.MergeContext(gate4.serverCtx, stream.Context())
	defer cancel()

	user, err := gate4.userStore.Find(ctx, req.UserId)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	client, err := gate4.clientPool.GetOrCreate(user)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return client.SubscribeQuotes(ctx, req.Symbol, func(q types.Quote) error {
		return stream.Send(&pb.Quote{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp,
			Ask:       q.Ask,
			Bid:       q.Bid,
		})
	})
}
