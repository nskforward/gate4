package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	transport "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Gate4Server struct {
	pb.UnimplementedAdminServer
	userStore     users.Store
	server        *transport.Server
	keychainStore *keychain.Store
}

func NewGate4Server(userStore users.Store, keychainStore *keychain.Store, opt ...transport.ServerOption) *Gate4Server {
	gate4 := &Gate4Server{
		server:        transport.NewServer(opt...),
		userStore:     userStore,
		keychainStore: keychainStore,
	}
	pb.RegisterAdminServer(gate4.server, gate4)
	return gate4
}

func (gate4 *Gate4Server) CreateCert(ctx context.Context, req *pb.CreateCertRequest) (*pb.CreateCertResponse, error) {
	if req.CommonName == "" {
		return nil, status.Error(codes.InvalidArgument, "common name cannot be empty")
	}

	certData, err := gate4.keychainStore.GenCert(req.CommonName, []byte(req.PrivateKey))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.CreateCertResponse{
		Cert: string(certData),
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

func (gate4 *Gate4Server) Run(ctx context.Context, network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("net.Listen error: %w", err)
	}

	errorc := make(chan error, 1)
	go func() {
		defer close(errorc)
		err := gate4.server.Serve(listener)
		if err != nil && err != transport.ErrServerStopped {
			errorc <- err
		}
	}()

	slog.Info("grpc server is ready to serve requests")

	select {
	case err := <-errorc:
		return err

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stopped := make(chan struct{})
		go func() {
			gate4.server.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-ctx.Done():
			gate4.server.Stop()
		}
		return nil
	}
}
