package broker

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/retries"
	"github.com/nskforward/gate4/pkg/types"
	"github.com/shopspring/decimal"
)

type positionStore struct {
	ctx     context.Context
	cancel  context.CancelFunc
	symbols map[string]*types.Position
	mx      sync.Mutex
}

func newPositionStore(client Client) (*positionStore, error) {
	ctx, cancel := context.WithCancel(context.Background())
	store := &positionStore{
		ctx:     ctx,
		cancel:  cancel,
		symbols: make(map[string]*types.Position),
	}
	t1 := time.Now()
	err := store.fill(client)
	if err != nil {
		return nil, err
	}
	go store.watch(client, t1)
	return store, err
}

func (store *positionStore) fill(client Client) error {
	reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	positions, err := client.GetPositions(reqCtx)
	if err != nil {
		return err
	}
	symbols := make(map[string]*types.Position)
	for _, pos := range positions {
		symbols[pos.Symbol] = &pos
	}
	store.symbols = symbols
	return nil
}

func (store *positionStore) watch(client Client, ignoreBefore time.Time) {
	slog.Debug("start watching account trades", "account_id", client.GetAccountInfo().AccountID)
	retry := store.initTradesRetry(client, ignoreBefore)
	err := retry.Do(store.ctx, func() error {
		return client.SubscribeAccountTrades(store.ctx, func(trade types.AccountTrade) bool {
			if trade.Timestamp < ignoreBefore.Unix() {
				// ignore old trades
				return true
			}
			store.calculate(trade)
			// always true
			return true
		})
	})
	if err != nil {
		slog.Error("stop watching account trades", "account_id", client.GetAccountInfo().AccountID, "error", err.Error())
		os.Exit(1)
	}
}

func (store *positionStore) initTradesRetry(client Client, ignoreBefore time.Time) *retries.Retry {
	return retries.New(retries.Config{
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
}

func (store *positionStore) calculate(trade types.AccountTrade) error {
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

	store.mx.Lock()
	defer store.mx.Unlock()

	pos, ok := store.symbols[trade.Symbol]
	if !ok {
		// no positions
		store.symbols[trade.Symbol] = &types.Position{
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
		delete(store.symbols, trade.Symbol)
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

/*

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

func (acc *accountPositions) GetAll() []types.Position {
	acc.mx.Lock()
	defer acc.mx.Unlock()

	positions := make([]types.Position, 0, len(acc.symbols))
	for _, pos := range acc.symbols {
		positions = append(positions, types.Position{
			Symbol: pos.Symbol,
			Price:  pos.Price,
			Size:   pos.Size,
		})
	}

	return positions
}
*/
