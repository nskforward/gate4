package transport

import (
	"context"
	"log/slog"

	"github.com/nskforward/gate4/internal/broker"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GatewayServer struct {
	pb.UnimplementedGatewayServer
	transport *grpcserv.GRPCServer
	broker    *broker.Broker
	logger    *slog.Logger
}

func NewGatewayServer(cfg config.Config, logger *slog.Logger, broker *broker.Broker) (*GatewayServer, error) {
	var server *grpcserv.GRPCServer

	if cfg.SSL.CA.CertPath != "" && cfg.SSL.Server.CertPath != "" && cfg.SSL.Server.KeyPath != "" {
		tlsConfig, err := grpcserv.MTLSConfig(cfg.SSL.CA.CertPath, cfg.SSL.Server.CertPath, cfg.SSL.Server.KeyPath)
		if err != nil {
			return nil, err
		}
		server = grpcserv.New(cfg.Gateway.ListenAddr, grpc.Creds(credentials.NewTLS(tlsConfig)))
		logger.Info("gateway mTLS enabled")
	} else {
		server = grpcserv.New(cfg.Gateway.ListenAddr)
		logger.Warn("gateway mTLS disabled")
	}

	s := &GatewayServer{
		transport: server,
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
	return s, nil
}

func (s *GatewayServer) Run(ctx context.Context) error {
	return s.transport.Run(ctx)
}

func (s *GatewayServer) GetAccountInfo(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	return s.broker.GetAccountInfo(ctx, in)
}
