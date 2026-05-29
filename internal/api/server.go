package api

import (
	"context"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/domain/users"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/ssl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedGate4Server
	userStorage   *users.UserStorage
	keychainStore *keychain.Store
	clientPool    *brokers.ClientPool
}

func NewServer(userStorage *users.UserStorage, keychainStore *keychain.Store, clientPool *brokers.ClientPool) *Server {
	return &Server{
		userStorage:   userStorage,
		keychainStore: keychainStore,
		clientPool:    clientPool,
	}
}

func (server *Server) Ping() string {
	return "pong"
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
	user, err := server.userStorage.Get(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return ConvertOutUser(user), nil
}

func (server *Server) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.ListUsersResponse, error) {
	users, err := server.userStorage.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ListUsersResponse{Users: ConvertOutUsers(users)}, nil
}

func (server *Server) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user := ConvertInUser(req)

	_, err := server.clientPool.GetOrCreate(&user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = server.userStorage.Create(ctx, &user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return ConvertOutUser(user), nil
}

func (server *Server) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.EmptyMessage, error) {
	user, err := server.userStorage.Get(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = server.clientPool.Delete(&user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = server.userStorage.Delete(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EmptyMessage{}, nil
}

func (server *Server) UpdateUser(ctx context.Context, req *pb.User) (*pb.EmptyMessage, error) {
	user := ConvertInUser(req)
	err := server.userStorage.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

/*
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
*/
