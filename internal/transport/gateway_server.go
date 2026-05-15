package transport

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nskforward/gate4/internal/broker"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GatewayServer struct {
	pb.UnimplementedGatewayServer
	serv   *grpcserv.GRPCServer
	broker *broker.Broker
}

func NewGatewayServer(cfg config.Config, broker *broker.Broker) (*GatewayServer, error) {
	var server *grpcserv.GRPCServer

	if cfg.SSL.CA.CertPath != "" && cfg.SSL.Server.CertPath != "" && cfg.SSL.Server.KeyPath != "" {
		tlsConfig, err := grpcserv.MTLSConfig(cfg.SSL.CA.CertPath, cfg.SSL.Server.CertPath, cfg.SSL.Server.KeyPath)
		if err != nil {
			return nil, err
		}
		server = grpcserv.New(cfg.Gateway.ListenAddr, grpc.Creds(credentials.NewTLS(tlsConfig)))
		slog.Info("gateway mTLS enabled")
	} else {
		server = grpcserv.New(cfg.Gateway.ListenAddr)
		slog.Warn("gateway mTLS disabled")
	}

	s := &GatewayServer{
		serv:   server,
		broker: broker,
	}
	pb.RegisterGatewayServer(s.serv, s)
	s.serv.OnListen = func() {
		slog.Info("gateway service started", "addr", cfg.Gateway.ListenAddr)
	}
	s.serv.OnStop = func() {
		slog.Info("gateway service stoppped")
	}
	return s, nil
}

func (s *GatewayServer) Run(ctx context.Context) error {
	return s.serv.Run(ctx)
}

/*
func (s *GatewayServer) GetPositions(ctx context.Context, req *pb.AccountRequest) (*pb.GetPositionsResponse, error) {
	account := s.broker.LookupAccount(req.AccountKey)
	if account == nil {
		return &pb.GetPositionsResponse{}, fmt.Errorf("unknown account")
	}

	positions, err := s.broker.GetPositions(ctx, account)
	if err != nil {
		return nil, err
	}

	return &pb.GetPositionsResponse{
		Positions: positions,
	}, nil
}
*/

func (s *GatewayServer) SubscribeQuoutes(req *pb.QuoteStreamRequest, serverStream pb.Gateway_SubscribeQuoutesServer) error {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "gateway subscribe quote stream",
		"id", id,
		"input", req.String(),
	)
	t1 := time.Now()
	err := s.broker.SubscribeQuoutes(serverStream.Context(), req.AccountKey, req.Symbol, func(q types.Quote) error {
		return serverStream.Send(convertQuote(q))
	})
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
		"error", err,
	)
	return err
}
