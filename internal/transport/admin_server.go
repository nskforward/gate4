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
)

type AdminServer struct {
	pb.UnimplementedAdminServer
	serv   *grpcserv.GRPCServer
	broker *broker.Broker
}

func NewAdminServer(cfg config.Config, broker *broker.Broker) (*AdminServer, error) {
	s := &AdminServer{
		serv:   grpcserv.New(cfg.Admin.ListenAddr),
		broker: broker,
	}
	pb.RegisterAdminServer(s.serv, s)
	s.serv.OnListen = func() {
		slog.Info("admin service started", "addr", cfg.Admin.ListenAddr)
	}
	s.serv.OnStop = func() {
		slog.Info("admin service stopped")
	}
	return s, nil
}

func (s *AdminServer) Run(ctx context.Context) error {
	go s.watch(ctx)
	return s.serv.Run(ctx)
}

func (s *AdminServer) ListAccounts(context.Context, *pb.EmptyMessage) (*pb.ListAccountsResponse, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin list acounts",
		"id", id,
	)
	t1 := time.Now()
	result := &pb.ListAccountsResponse{
		Items: broker.ExportAccounts(s.broker.Accounts()),
	}
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
	)
	return result, nil
}

func (s *AdminServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.EmptyMessage, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin add acount",
		"id", id,
		"broker", req.Account.BrokerId,
		"account_id", req.Account.Id,
		"valid_until", req.Account.ValidUntil,
	)
	t1 := time.Now()
	account := broker.ImportAccount(req.Account)
	err := s.broker.AddAccount(account)
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
	)
	return &pb.EmptyMessage{}, err
}

func (s *AdminServer) DeleteAccount(_ context.Context, req *pb.AccountRequest) (*pb.EmptyMessage, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin del acount",
		"id", id,
		"input", req.String(),
	)
	t1 := time.Now()
	err := s.broker.DelAccount(req.AccountKey)
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
	)
	return &pb.EmptyMessage{}, err
}

func (s *AdminServer) SubscribeQuoutes(req *pb.SymbolRequest, serverStream pb.Admin_SubscribeQuoutesServer) error {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin subscribe quote stream",
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

func (s *AdminServer) SubscribeAccountTrades(req *pb.AccountRequest, serverStream pb.Admin_SubscribeAccountTradesServer) error {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin subscribe account trades",
		"id", id,
		"input", req.String(),
	)
	t1 := time.Now()
	err := s.broker.SubscribeAccountTrades(serverStream.Context(), req.AccountKey, func(t types.AccountTrade) error {
		return serverStream.Send(convertAccountTrade(t))
	})
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
		"error", err,
	)
	return err
}

func (s *AdminServer) GetPositions(ctx context.Context, req *pb.AccountRequest) (*pb.ListPositions, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin get positions",
		"id", id,
		"input", req.String(),
	)
	positions, err := s.broker.GetPositions(ctx, req.AccountKey)
	t1 := time.Now()
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
		"error", err,
	)
	return convertPositions(positions), nil
}

func (s *AdminServer) GetPosition(ctx context.Context, req *pb.SymbolRequest) (*pb.Position, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin get position",
		"id", id,
		"input", req.String(),
	)
	info, err := s.broker.GetPosition(ctx, req.AccountKey, req.Symbol)
	t1 := time.Now()
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
		"error", err,
	)
	return convertPosition(info), nil
}

func (s *AdminServer) GetSchedule(ctx context.Context, req *pb.SymbolRequest) (*pb.GetScheduleResponse, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin get shedule",
		"id", id,
		"input", req.String(),
	)
	sessions, err := s.broker.GetSchedule(ctx, req.AccountKey, req.Symbol)
	t1 := time.Now()
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
		"error", err,
	)
	if err != nil {
		return nil, err
	}
	return convertSessions(sessions), nil
}

func (s *AdminServer) GetAsset(ctx context.Context, req *pb.SymbolRequest) (*pb.GetAssetResponse, error) {
	id := uuid.NewString()
	slog.Debug("start call",
		"name", "admin get asset",
		"id", id,
		"input", req.String(),
	)
	asset, err := s.broker.GetAsset(ctx, req.AccountKey, req.Symbol)
	t1 := time.Now()
	slog.Debug("finish call",
		"id", id,
		"duration", time.Since(t1).String(),
		"error", err,
	)
	if err != nil {
		return nil, err
	}
	return convertAsset(asset), nil
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

	slog.Warn("account will expire soon", "broker_id", account.Broker, "account_id", account.ID, "expires_in_days", days)
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
