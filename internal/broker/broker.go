package broker

import (
	"context"
	"fmt"
	"sync"

	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/types"
)

type Broker struct {
	accounts     *AccountStore
	finamClients *finam.MultiClient
	mx           sync.Mutex
}

func NewBroker() (*Broker, error) {
	accounts, err := NewAccountStore("data/accounts.json")
	if err != nil {
		return nil, err
	}
	return &Broker{
		accounts:     accounts,
		finamClients: finam.NewMultiClient(),
	}, nil
}

func (b *Broker) AddAccount(account *Account) error {
	return b.accounts.Set(account.Key(), account)
}

func (b *Broker) DelAccount(key string) error {
	return b.accounts.Del(key)
}

func (b *Broker) Accounts() []*Account {
	return b.accounts.Accounts()
}

func (b *Broker) SubscribeQuoutes(ctx context.Context, req *pb.QuoteStreamRequest, send func(*pb.QuoteStreamResponse) error) error {
	client, err := b.getClient(req.AccountKey)
	if err != nil {
		return err
	}

	stream := client.SubscribeQuotes(ctx, req.Symbol)
	defer stream.Close()

	for q := range stream.Range() {
		err := send(&pb.QuoteStreamResponse{
			Symbol:    q.Symbol,
			Timestamp: q.Timestamp,
			Ask:       q.Ask.Price,
			Bid:       q.Bid.Price,
		})
		if err != nil {
			return err
		}
	}
	return stream.Err()
}

func (b *Broker) GetPositions(ctx context.Context, req *pb.AccountRequest) (*pb.GetPositionsResponse, error) {
	client, err := b.getClient(req.AccountKey)
	if err != nil {
		return nil, err
	}
	info, err := client.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	positions := make([]*pb.Position, 0, len(info.Positions))
	for _, pos := range info.Positions {
		positions = append(positions, &pb.Position{
			Symbol: pos.Symbol,
			Price:  pos.Price,
			Size:   pos.Size,
			Profit: pos.Profit,
		})
	}
	return &pb.GetPositionsResponse{
		Positions: positions,
	}, nil
}

func (b *Broker) GetSchedule(ctx context.Context, req *pb.GetScheduleRequest) (*pb.GetScheduleResponse, error) {
	client, err := b.getClient(req.AccountKey)
	if err != nil {
		return nil, err
	}
	sessions, current, err := client.GetSchedule(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	return &pb.GetScheduleResponse{
		CurrentSession: toPbSession(current),
		Sessions:       toPbSessions(sessions),
	}, nil
}

func (b *Broker) getClient(key string) (types.BrokerClient, error) {
	account, ok := b.accounts.Lookup(key)
	if !ok {
		return nil, fmt.Errorf("unknown account key: %s", key)
	}
	switch account.Broker {
	case "finam":
		return b.finamClients.Get(account.ID, account.Secret)
	}
	return nil, fmt.Errorf("client not registered for broker: %s", account.Broker)
}
