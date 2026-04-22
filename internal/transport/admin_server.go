package transport

import (
	"context"
	"log/slog"

	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
)

type AdminServer struct {
	pb.UnimplementedAdminServer
	transport *grpcserv.GRPCServer
	logger    *slog.Logger
}

func NewAdminServer(cfg config.Config, logger *slog.Logger) *AdminServer {
	s := &AdminServer{
		transport: grpcserv.New(cfg.Admin.ListenAddr),
		logger:    logger,
	}
	pb.RegisterAdminServer(s.transport, s)
	s.transport.OnListen = func() {
		logger.Info("admin service started", "addr", cfg.Admin.ListenAddr)
	}
	s.transport.OnStop = func() {
		logger.Info("admin service stoppped")
	}
	return s
}

func (s *AdminServer) Run(ctx context.Context) error {
	return s.transport.Run(ctx)
}
