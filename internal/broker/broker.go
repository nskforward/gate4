package broker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/finam"
	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/streams"
	"github.com/nskforward/gate4/pkg/types"
)

type Broker struct {
	accounts            *AccountStore
	finamClients        *finam.MultiClient
	mx                  sync.Mutex
	quoteStreams        *streams.Store[types.Quote]
	accountTradeStreams *streams.Store[types.AccountTrade]
	scheduleCache       *ScheduleCache
	positionStore       *PositionStore
}

func NewBroker() (*Broker, error) {
	accounts, err := NewAccountStore("data/accounts.json")
	if err != nil {
		return nil, err
	}
	return &Broker{
		accounts:            accounts,
		finamClients:        finam.NewMultiClient(),
		quoteStreams:        streams.NewStore[types.Quote](context.Background(), 1),
		accountTradeStreams: streams.NewStore[types.AccountTrade](context.Background(), 32),
		scheduleCache:       NewScheduleCache(),
		positionStore:       NewPositionStore(),
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

func (b *Broker) SubscribeQuoutes(ctx context.Context, accountKey, symbol string, send func(types.Quote) error) error {
	client, err := b.getClient(accountKey)
	if err != nil {
		return err
	}
	stream := b.quoteStreams.Subscribe(ctx, symbol, func(topicCtx context.Context, publish func(data types.Quote) bool) error {
		retry := retries.New(retries.Config{
			InitialDelay:  500 * time.Millisecond,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			MaxAttempts:   10,
			MaxJitter:     time.Second,
			OnError: func(err error) {
				slog.Error("unexpectedly quote stream closed", "broker", client.GetAccountInfo().BrokerID, "account_id", client.GetAccountInfo().AccountID, "symbol", symbol, "error", err.Error())
			},
			OnAttempt: func(attempt int) {
				slog.Debug("try to connect to quote stream", "broker", client.GetAccountInfo().BrokerID, "account_id", client.GetAccountInfo().AccountID, "symbol", symbol, "attempt", attempt)
			},
		})
		return retry.Do(topicCtx, func() error {
			return client.SubscribeQuotes(topicCtx, symbol, publish)
		})
	})
	for q := range stream.Range() {
		err := send(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) SubscribeAccountTrades(ctx context.Context, accountKey string, send func(types.AccountTrade) error) error {
	client, err := b.getClient(accountKey)
	if err != nil {
		return err
	}
	stream := b.accountTradeStreams.Subscribe(ctx, accountKey, func(topicCtx context.Context, publish func(data types.AccountTrade) bool) error {
		retry := retries.New(retries.Config{
			InitialDelay:  500 * time.Millisecond,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			MaxAttempts:   10,
			MaxJitter:     time.Second,
		})
		return retry.Do(topicCtx, func() error {
			return client.SubscribeAccountTrades(topicCtx, publish)
		})
	})
	for q := range stream.Range() {
		err := send(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) GetAccountInfo(ctx context.Context, accountKey string) (types.AccountInfo, error) {
	client, err := b.getClient(accountKey)
	if err != nil {
		return types.AccountInfo{}, err
	}
	return client.GetAccountInfo(), nil
}

func (b *Broker) GetPositions(ctx context.Context, accountKey string) ([]types.Position, error) {
	client, err := b.getClient(accountKey)
	if err != nil {
		return nil, err
	}
	acc, err := b.positionStore.Get(ctx, client)
	if err != nil {
		return nil, err
	}
	return acc.GetAll(), nil
}

func (b *Broker) GetPosition(ctx context.Context, accountKey, symbol string) (types.Position, error) {
	client, err := b.getClient(accountKey)
	if err != nil {
		return types.Position{}, err
	}
	acc, err := b.positionStore.Get(ctx, client)
	if err != nil {
		return types.Position{}, err
	}
	return acc.Get(symbol), nil
}

func (b *Broker) GetSchedule(ctx context.Context, accountKey, symbol string) ([]types.Session, error) {
	client, err := b.getClient(accountKey)
	if err != nil {
		return nil, err
	}
	return b.scheduleCache.GetSessions(ctx, client, symbol)
}

func (b *Broker) GetAsset(ctx context.Context, accountKey, symbol string) (types.AssetInfo, error) {
	client, err := b.getClient(accountKey)
	if err != nil {
		return types.AssetInfo{}, err
	}
	return client.GetAsset(ctx, symbol)
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
