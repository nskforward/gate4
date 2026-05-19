package broker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/types"
	"github.com/shopspring/decimal"
)

type PositionStore struct {
	accounts map[Client]*accountPositions
	mx       sync.Mutex
}

type accountPositions struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols map[string]*types.Position
	mx      sync.Mutex
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		accounts: make(map[Client]*accountPositions),
	}
}

func (store *PositionStore) Get(ctx context.Context, client Client) (*accountPositions, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	positions, ok := store.accounts[client]
	if ok {
		return positions, nil
	}

	t1 := time.Now()

	allPositions, err := client.GetPositions(ctx)
	if err != nil {
		return nil, err
	}

	acc := newAccountPositions(allPositions)

	store.accounts[client] = acc

	go store.watchTrades(acc, client, t1)

	return acc, nil
}

func (store *PositionStore) Del(client Client) {
	store.mx.Lock()
	defer store.mx.Unlock()

	acc, ok := store.accounts[client]
	if !ok {
		return
	}

	acc.cancel()
	delete(store.accounts, client)
}

func (store *PositionStore) watchTrades(acc *accountPositions, client Client, ignoreBefore time.Time) {
	slog.Debug("start watching account trades", "account_id", client.GetAccountInfo().AccountID)

	retry := retries.New(retries.Config{
		InitialDelay:  time.Second,
		MaxDelay:      time.Minute,
		BackoffFactor: 2,
		MaxAttempts:   0,
		MaxJitter:     time.Second,
		OnError: func(err error) {
			slog.Error("unexpectedly account trades stream closed", "broker", client.GetAccountInfo().BrokerID, "account_id", client.GetAccountInfo().AccountID, "error", err.Error())
		},
		OnAttempt: func(attempt int) {
			slog.Debug("try to connect to account trades stream", "broker", client.GetAccountInfo().BrokerID, "account_id", client.GetAccountInfo().AccountID, "attempt", attempt)
		},
	})

	err := retry.Do(acc.ctx, func() error {
		return client.SubscribeAccountTrades(acc.ctx, func(trade types.AccountTrade) bool {
			if trade.Timestamp < ignoreBefore.Unix() {
				// ignore old trades
				return true
			}

			acc.calculate(trade)

			// always true
			return true
		})
	})
	if err != nil {
		slog.Debug("stop watching account trades", "account_id", client.GetAccountInfo().AccountID, "error", err.Error())
		panic(err)
	}
}

func newAccountPositions(in []types.Position) *accountPositions {
	symbols := make(map[string]*types.Position)
	for _, pos := range in {
		symbols[pos.Symbol] = &pos
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &accountPositions{
		ctx:     ctx,
		cancel:  cancel,
		symbols: symbols,
	}

}

func (acc *accountPositions) Get(symbol string) types.Position {
	acc.mx.Lock()
	defer acc.mx.Unlock()
	pos, ok := acc.symbols[symbol]
	if !ok {
		return types.Position{}
	}
	return types.Position{
		Symbol: pos.Symbol,
		Price:  pos.Price,
		Size:   pos.Size,
	}
}

func (acc *accountPositions) calculate(trade types.AccountTrade) error {
	tradeSize, err := decimal.NewFromString(trade.Size)
	if err != nil {
		return err
	}

	tradePrice, err := decimal.NewFromString(trade.Price)
	if err != nil {
		return err
	}

	if tradeSize.Sign() == 0 {
		return nil
	}

	acc.mx.Lock()
	defer acc.mx.Unlock()

	pos, ok := acc.symbols[trade.Symbol]
	if !ok {
		// no positions
		acc.symbols[trade.Symbol] = &types.Position{
			Symbol: trade.Symbol,
			Price:  trade.Price,
			Size:   trade.Size,
		}
		return nil
	}

	posSize, err := decimal.NewFromString(pos.Size)
	if err != nil {
		return err
	}

	posPrice, err := decimal.NewFromString(pos.Price)
	if err != nil {
		return err
	}

	if tradeSize.Sign() == posSize.Sign() {
		// add more in the same direction
		totalSize := posSize.Add(tradeSize)
		pos.Price = posSize.Mul(posPrice).Add(tradeSize.Mul(tradePrice)).Div(totalSize).StringFixedBank(2)
		pos.Size = totalSize.String()
		return nil
	}

	totalSize := posSize.Add(tradeSize)
	if totalSize.IsZero() {
		// full closure
		pos.Size = ""
		delete(acc.symbols, trade.Symbol)
		return nil
	}

	cmp := tradeSize.Abs().Cmp(posSize.Abs())
	if cmp < 0 {
		// partial closure
		pos.Size = totalSize.String()
		return nil
	}

	// coup
	pos.Price = trade.Price
	pos.Size = totalSize.String()
	return nil
}
