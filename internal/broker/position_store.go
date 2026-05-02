package broker

import (
	"log/slog"
	"sync"

	"github.com/nskforward/gate4/pkg/pb"
)

type PositionStore struct {
	items map[string][]*pb.Position
	mx    sync.RWMutex
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		items: make(map[string][]*pb.Position),
	}
}

func (s *PositionStore) Update(account *Account, positions []*pb.Position) {
	for _, pos := range positions {
		slog.Debug("save position", "symbol", pos.Symbol, "price", pos.AveragePrice, "size", pos.Size)
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	if positions == nil {
		delete(s.items, account.Key())
		return
	}

	s.items[account.Key()] = positions
}

func (s *PositionStore) Get(account *Account) []*pb.Position {
	s.mx.RLock()
	defer s.mx.RUnlock()
	positions, ok := s.items[account.Key()]
	if !ok {
		return []*pb.Position{}
	}
	return positions
}
