package transport

import (
	"context"
	"log/slog"
	"time"

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

func NewAdminServer(cfg config.Config, logger *slog.Logger, accountStore *store.AccountStore) *AdminServer {
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
	return s
}

func (s *AdminServer) Run(ctx context.Context) error {
	go s.watch(ctx)
	return s.transport.Run(ctx)
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
	s.accountStore.Del(req.Key)
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
	items := s.accountStore.List()
	for _, item := range items {
		t := time.Unix(item.ValidUntil, 0)
		p := t.Sub(now)
		s.notifyAccountExpiration(item, p.Hours())
	}
}

func (s *AdminServer) notifyAccountExpiration(account *pb.Account, hours float64) {
	if hours > 24*7 {
		// greater than a week
		return
	}

	days := int(hours / 24)

	s.logger.Warn("account will expire soon", "broker_id", account.BrokerId, "account_id", account.Id, "expires_in_days", days)
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
