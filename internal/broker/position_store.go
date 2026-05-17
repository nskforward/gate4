package broker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/types"
)

type PositionStore struct {
	accounts map[Client][]types.Position
	mx       sync.Mutex
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		accounts: make(map[Client][]types.Position),
	}
}

func (store *PositionStore) Get(ctx context.Context, client Client) ([]types.Position, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	positions, ok := store.accounts[client]
	if ok {
		return positions, nil
	}

	info, err := client.GetAccount(ctx)
	if err != nil {
		return nil, err
	}

	store.accounts[client] = info.Positions

	go store.watchTrades(ctx, client)

	return info.Positions, nil
}

func (store *PositionStore) watchTrades(ctx context.Context, client Client) {
	slog.Debug("start watching account trades", "account_id", client.GetAccountID())

	retry := retries.New(retries.Config{
		InitialDelay:  time.Second,
		MaxDelay:      time.Minute,
		BackoffFactor: 2,
		MaxAttempts:   0,
		MaxJitter:     time.Second,
		OnError: func(err error) {
			slog.Error("unexpectedly account trades stream closed", "broker", client.GetBrokerID(), "account_id", client.GetAccountID(), "error", err.Error())
		},
		OnAttempt: func(attempt int) {
			slog.Debug("try to connect to account trades stream", "broker", client.GetBrokerID(), "account_id", client.GetAccountID(), "attempt", attempt)
		},
	})

	err := retry.Do(ctx, func() error {
		return client.SubscribeAccountTrades(ctx, func(trade types.AccountTrade) bool {
			store.mx.Lock()
			defer store.mx.Unlock()

			// TODO !!!
			fmt.Println("account trade:", trade)

			return true
		})
	})
	if err != nil {
		slog.Debug("stop watching account trades", "account_id", client.GetAccountID(), "error", err.Error())
		panic(err)
	}
}
