package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/streams"
	"github.com/nskforward/gate4/pkg/types"
)

type Broker struct {
	accounts      *AccountStore
	finamClients  *finam.MultiClient
	mx            sync.Mutex
	quoteStreams  *streams.Store[types.Quote]
	scheduleCache *ScheduleCache
}

func NewBroker() (*Broker, error) {
	accounts, err := NewAccountStore("data/accounts.json")
	if err != nil {
		return nil, err
	}
	return &Broker{
		accounts:      accounts,
		finamClients:  finam.NewMultiClient(),
		quoteStreams:  streams.NewStore[types.Quote](context.Background(), 1),
		scheduleCache: NewScheduleCache(),
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

	stream := b.quoteStreams.Subscribe(ctx, req.Symbol, func(topicCtx context.Context, publish func(data types.Quote) bool) error {
		retry := retries.New(retries.Config{
			InitialDelay:  500 * time.Millisecond,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			MaxAttempts:   10,
			MaxJitter:     time.Second,
		})
		return retry.Do(topicCtx, func() error {
			return client.SubscribeQuotes(topicCtx, req.Symbol, publish)
		})
	})

	for q := range stream.Range() {
		err := send(toPbQuote(q))
		if err != nil {
			return err
		}
	}
	return nil
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
	sessions, err := b.scheduleCache.GetSessions(ctx, client, req.Symbol)
	if err != nil {
		return nil, err
	}
	return &pb.GetScheduleResponse{
		Sessions: toPbSessions(sessions),
	}, nil
}

func (b *Broker) getClient(key string) (Client, error) {
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
