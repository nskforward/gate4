package transport

import (
	"context"
	"log/slog"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/store"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
)

type AdminServer struct {
	pb.UnimplementedAdminServer
	transport    *grpcserv.GRPCServer
	logger       *slog.Logger
	accountStore *store.AccountStore
}

func NewAdminServer(cfg config.Config, logger *slog.Logger) (*AdminServer, error) {
	accountStore, err := store.NewAccountStore(store.NewAccountFileProvider(cfg.StoreDir))
	if err != nil {
		return nil, err
	}

	s := &AdminServer{
		transport:    grpcserv.New(cfg.Admin.ListenAddr),
		logger:       logger,
		accountStore: accountStore,
	}
	pb.RegisterAdminServer(s.transport, s)
	s.transport.OnListen = func() {
		logger.Info("admin service started", "addr", cfg.Admin.ListenAddr)
	}
	s.transport.OnStop = func() {
		logger.Info("admin service stoppped")
	}
	return s, nil
}

func (s *AdminServer) Run(ctx context.Context) error {
	err := s.transport.Run(ctx)
	if err != nil {
		return err
	}
	return s.accountStore.Close()
}

func (s *AdminServer) ListAccounts(context.Context, *pb.EmptyMessage) (*pb.ListAccountsResponse, error) {
	return &pb.ListAccountsResponse{
		Items: s.accountStore.List(),
	}, nil
}

func (s *AdminServer) AddAccount(_ context.Context, req *pb.AddAccountRequest) (*pb.EmptyMessage, error) {
	err := s.accountStore.Set(req.Account)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *AdminServer) DeleteAccount(_ context.Context, req *pb.DeleteAccountRequest) (*pb.EmptyMessage, error) {
	s.accountStore.Del(req.Id)
	return &pb.EmptyMessage{}, nil
}
