package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
	transport "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminServer struct {
	pb.UnimplementedAdminServer
	addr          string
	userStore     users.Store
	transport     *transport.Server
	keychainStore *keychain.Store
}

func NewAdminServer(ctx context.Context, cfg config.Config) (*AdminServer, error) {
	keychainStore, err := keychain.NewStore(ctx, cfg)
	if err != nil {
		return nil, err
	}

	admin := &AdminServer{
		addr:          cfg.Admin.ListenAddr,
		userStore:     users.NewFileStorage(),
		transport:     transport.NewServer(),
		keychainStore: keychainStore,
	}
	pb.RegisterAdminServer(admin.transport, admin)
	return admin, nil
}

func (admin *AdminServer) FindUser(ctx context.Context, req *pb.UserID) (*pb.User, error) {
	user, err := admin.userStore.Find(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return convertOutUser(user), nil
}

func (admin *AdminServer) ListUsers(ctx context.Context, req *pb.EmptyMessage) (*pb.ListUsersResponse, error) {
	users, err := admin.userStore.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ListUsersResponse{Users: convertOutUsers(users)}, nil
}

func (admin *AdminServer) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user := convertInUser(req)
	err := admin.userStore.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertOutUser(user), nil
}

func (admin *AdminServer) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.EmptyMessage, error) {
	err := admin.userStore.Delete(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EmptyMessage{}, nil
}

func (admin *AdminServer) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.EmptyMessage, error) {
	err := admin.userStore.Block(ctx, req.UserId, req.Blocked)
	return &pb.EmptyMessage{}, err
}

func (admin *AdminServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := admin.userStore.Update(ctx, req.UserId, req.Secret, time.Unix(req.Expires, 0))
	if err != nil {
		return nil, err
	}
	return convertOutUser(user), nil
}

func (admin *AdminServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", admin.addr)
	if err != nil {
		return fmt.Errorf("net.Listen error: %w", err)
	}

	errorc := make(chan error, 1)
	go func() {
		defer close(errorc)
		err := admin.transport.Serve(listener)
		if err != nil && err != transport.ErrServerStopped {
			errorc <- err
		}
	}()

	slog.Info("admin server is ready to serve requests")

	select {
	case err := <-errorc:
		return err

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stopped := make(chan struct{})
		go func() {
			admin.transport.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-ctx.Done():
			admin.transport.Stop()
		}
		return nil
	}
}
