package transport

import (
	"context"
	"log/slog"

	"github.com/nskforward/gate4/internal/broker"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
)

type GatewayServer struct {
	pb.UnimplementedGatewayServer
	transport *grpcserv.GRPCServer
	broker    *broker.Broker
	logger    *slog.Logger
}

func NewGatewayServer(cfg config.Config, logger *slog.Logger, broker *broker.Broker) *GatewayServer {
	s := &GatewayServer{
		transport: grpcserv.New(cfg.Gateway.ListenAddr),
		broker:    broker,
		logger:    logger,
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
	return s.broker.GetAccount(ctx, in)
}
