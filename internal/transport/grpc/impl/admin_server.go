package impl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nskforward/gate4/internal/broker"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/pkg/grpcserv"
	"github.com/nskforward/gate4/pkg/pb"
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

func (s *AdminServer) QuoteStream(req *pb.QuoteStreamRequest, serverStream pb.Admin_QuoteStreamServer) error {
	client, err := s.broker.Client(req.AccountKey)
	if err != nil {
		return err
	}
	stream, err := client.SubscribeQuotes(serverStream.Context(), []string{req.Symbol})
	if err != nil {
		return err
	}
	for q := range stream.Range() {
		err := serverStream.Send(&pb.QuoteStreamResponse{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp,
			Ask:       q.Ask.Price,
			Bid:       q.Bid.Price,
		})
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("client disconnected")
}

func (s *AdminServer) ListAccounts(context.Context, *pb.EmptyMessage) (*pb.ListAccountsResponse, error) {
	return &pb.ListAccountsResponse{
		Items: broker.ExportAccounts(s.broker.Accounts()),
	}, nil
}

func (s *AdminServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.EmptyMessage, error) {
	account := broker.ImportAccount(req.Account)
	err := s.broker.AddAccount(account)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

func (s *AdminServer) DeleteAccount(_ context.Context, req *pb.AccountRequest) (*pb.EmptyMessage, error) {
	err := s.broker.DelAccount(req.AccountKey)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyMessage{}, nil
}

/*
	func (s *AdminServer) GetPositions(ctx context.Context, req *pb.AccountRequest) (*pb.GetPositionsResponse, error) {
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

	func (s *AdminServer) GetSchedule(ctx context.Context, req *pb.GetScheduleRequest) (*pb.GetScheduleResponse, error) {
		account := s.broker.LookupAccount(req.AccountKey)
		if account == nil {
			return nil, fmt.Errorf("unknown account")
		}
		sessions, current, err := s.broker.GetSchedule(ctx, account, req.Symbol)
		if err != nil {
			return nil, err
		}
		return &pb.GetScheduleResponse{
			CurrentSession: current,
			Sessions:       sessions,
		}, nil
	}
*/

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
