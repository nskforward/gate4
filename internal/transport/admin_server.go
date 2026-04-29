package transport

import (
	"context"
	"log/slog"
	"time"

	"github.com/nskforward/gate4/internal/broker"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AdminServer struct {
	pb.UnimplementedAdminServer
	transport *grpcserv.GRPCServer
	logger    *slog.Logger
	broker    *broker.Broker
}

func NewAdminServer(cfg config.Config, logger *slog.Logger, broker *broker.Broker) (*AdminServer, error) {

	var server *grpcserv.GRPCServer

	if cfg.SSL.CA.CertPath != "" && cfg.SSL.Server.CertPath != "" && cfg.SSL.Server.KeyPath != "" {
		tlsConfig, err := grpcserv.MTLSConfig(cfg.SSL.CA.CertPath, cfg.SSL.Server.CertPath, cfg.SSL.Server.KeyPath)
		if err != nil {
			return nil, err
		}
		server = grpcserv.New(cfg.Admin.ListenAddr, grpc.Creds(credentials.NewTLS(tlsConfig)))
		logger.Info("mTLS enabled")
	} else {
		server = grpcserv.New(cfg.Admin.ListenAddr)
		logger.Warn("mTLS disabled")
	}

	s := &AdminServer{
		transport: server,
		logger:    logger,
		broker:    broker,
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
	go s.watch(ctx)
	return s.transport.Run(ctx)
}

func (s *AdminServer) ListAccounts(context.Context, *pb.EmptyMessage) (*pb.ListAccountsResponse, error) {
	return &pb.ListAccountsResponse{
		Items: broker.ExportAccounts(s.broker.Accounts()),
	}, nil
}

func (s *AdminServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.EmptyMessage, error) {
	account := broker.ImportAccount(req.Account)

	client, err := s.broker.LookupClient(account)
	if err != nil {
		return nil, err
	}

	_, err = client.GetAccountInfo(ctx, account)
	if err != nil {
		return nil, err
	}

	err = s.broker.AddAccount(account)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *AdminServer) DeleteAccount(_ context.Context, req *pb.DeleteAccountRequest) (*pb.EmptyMessage, error) {
	err := s.broker.DeleteAccount(req.Key)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *AdminServer) watch(ctx context.Context) {
	for {
		sleep := getSleep()
		select {
		case <-ctx.Done():
			return
		case <-time.After(sleep):
			s.checkAccounts()
		}
	}
}

func (s *AdminServer) checkAccounts() {
	now := time.Now()
	items := s.broker.Accounts()
	for _, item := range items {
		p := item.ValidUntil.Sub(now)
		s.notifyAccountExpiration(item, p.Hours())
	}
}

func (s *AdminServer) notifyAccountExpiration(account *broker.Account, hours float64) {
	if hours > 24*7 {
		// greater than a week
		return
	}

	days := int(hours / 24)

	s.logger.Warn("account will expire soon", "broker_id", account.Broker, "account_id", account.ID, "expires_in_days", days)
}

func getSleep() time.Duration {
	now := time.Now()
	if now.Hour() < 9 {
		targetDate := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local)
		return targetDate.Sub(now)
	}
	targetDate := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, time.Local)
	return targetDate.Sub(now)
}
