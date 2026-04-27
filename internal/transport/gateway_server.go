package transport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/store"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
)

type GatewayServer struct {
	pb.UnimplementedGatewayServer
	transport    *grpcserv.GRPCServer
	accountStore *store.AccountStore
	logger       *slog.Logger
}

func NewGatewayServer(cfg config.Config, logger *slog.Logger, accountStore *store.AccountStore) *GatewayServer {
	s := &GatewayServer{
		transport:    grpcserv.New(cfg.Gateway.ListenAddr),
		accountStore: accountStore,
		logger:       logger,
	}
	pb.RegisterGatewayServer(s.transport, s)
	s.transport.OnListen = func() {
		logger.Info("gateway service started", "addr", cfg.Gateway.ListenAddr)
	}
	s.transport.OnStop = func() {
		logger.Info("gateway service stoppped")
	}
	return s
}

func (s *GatewayServer) Run(ctx context.Context) error {
	return s.transport.Run(ctx)
}

func (s *GatewayServer) GetAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	account := s.accountStore.Get(in.AccountKey)
	if account == nil {
		return nil, errors.New("unknown account key")
	}

	switch account.BrokerId {
	case "finam":
		// TODO implement finam client store: pkg/finam/store.go
		// client := finam.GetClient(account)
		client, err := s.finamAccounts.Get(in.AccountId)
		if err != nil {
			return nil, fmt.Errorf("finam client: %w", err)
		}
		resp, err := client.GetAccountInfo(ctx, in.AccountId)
		if err != nil {
			return nil, fmt.Errorf("finam communication error: %w", err)
		}
		return &pb.AccountResponse{
			BrokerId:  "finam",
			AccountId: resp.AccountId,
		}, nil
	}
	return nil, fmt.Errorf("broker '%s' not supported", in.BrokerId)
}
