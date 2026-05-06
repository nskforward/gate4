package broker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/nskforward/gate4/internal/broker/types"
	"github.com/nskforward/gate4/pkg/finam"
)

type PositionStore struct {
	accounts map[string]positionStoreLeaf
	mx       sync.Mutex
}

type positionStoreLeaf struct {
	Timestamp int64
	Positions []types.Position
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		accounts: make(map[string]positionStoreLeaf),
	}
}

func (s *PositionStore) Get(ctx context.Context, client *finam.Client, accountID string) ([]types.Position, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	leaf, ok := s.accounts[accountID]

	cache := "hit"
	defer func() {
		slog.Debug("get positions", "broker", "finam", "account", accountID, "cache", cache)
	}()

	if ok && time.Now().Unix()-leaf.Timestamp < 5 {
		return leaf.Positions, nil
	}

	cache = "miss"

	resp, err := client.GetAccountInfo(ctx, accountID)
	if err != nil {
		return nil, err
	}

	positions := make([]types.Position, 0, len(resp.Positions))
	for _, pos := range resp.Positions {
		positions = append(positions, types.Position{
			Symbol: pos.Symbol,
			Price:  pos.AveragePrice.Value,
			Size:   pos.Quantity.Value,
			Profit: pos.UnrealizedPnl.Value,
		})
	}

	leaf.Timestamp = time.Now().Unix()
	leaf.Positions = positions

	s.accounts[accountID] = leaf

	return positions, nil
}
