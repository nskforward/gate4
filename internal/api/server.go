package api

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/ssl"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedGate4Server
	userStore     users.Store
	keychainStore *keychain.Store
	clientPool    *brokers.ClientPool
}

func NewServer(userStore users.Store, keychainStore *keychain.Store, clientPool *brokers.ClientPool) *Server {
	return &Server{
		userStore:     userStore,
		keychainStore: keychainStore,
		clientPool:    clientPool,
	}
}

func (server *Server) CreateCert(ctx context.Context, req *pb.CreateCertRequest) (*pb.CreateCertResponse, error) {
	if req.CommonName == "" {
		return nil, status.Error(codes.InvalidArgument, "common name cannot be empty")
	}
	key, err := ssl.ParsePrivateKey([]byte(req.PrivateKey))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	cert, err := server.keychainStore.Generate(req.CommonName, key)
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

func (server *Server) FindUser(ctx context.Context, req *pb.UserID) (*pb.User, error) {
	user, err := server.userStore.Find(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return ConvertOutUser(user), nil
}

func (server *Server) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.ListUsersResponse, error) {
	console.LogDebug("api call method ListUsers")
	users, err := server.userStore.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ListUsersResponse{Users: ConvertOutUsers(users)}, nil
}

func (server *Server) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user := ConvertInUser(req)

	_, err := server.clientPool.GetOrCreate(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = server.userStore.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return ConvertOutUser(user), nil
}

func (server *Server) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.EmptyMessage, error) {
	user, err := server.userStore.Find(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = server.clientPool.Delete(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = server.userStore.Delete(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EmptyMessage{}, nil
}

func (server *Server) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.EmptyMessage, error) {
	err := server.userStore.Block(ctx, req.UserId, req.Blocked)
	return &pb.EmptyMessage{}, err
}

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := server.userStore.Update(ctx, req.UserId, req.Secret, time.Unix(req.Expires, 0))
	if err != nil {
		return nil, err
	}
	return ConvertOutUser(user), nil
}

func (server *Server) SubscribeQuotes(req *pb.SymbolRequest, stream grpc.ServerStreamingServer[pb.Quote]) error {
	user, err := server.userStore.Find(stream.Context(), req.UserId)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	client, err := server.clientPool.GetOrCreate(user)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return client.SubscribeQuotes(stream.Context(), req.Symbol, func(q types.Quote) error {
		return stream.Send(&pb.Quote{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp,
			Ask:       q.Ask,
			Bid:       q.Bid,
		})
	})
}
