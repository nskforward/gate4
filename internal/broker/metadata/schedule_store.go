package metadata

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/types"
)

type ScheduleStore struct {
	client  types.Client
	ctx     context.Context
	cancel  context.CancelFunc
	symbols map[string][]types.Session
	mx      sync.Mutex
}

func NewScheduleStore(client types.Client) *ScheduleStore {
	ctx, cancel := context.WithCancel(context.Background())
	return &ScheduleStore{
		ctx:     ctx,
		cancel:  cancel,
		client:  client,
		symbols: make(map[string][]types.Session),
	}
}

func (store *ScheduleStore) Get(symbol string) ([]types.Session, error) {
	store.mx.Lock()
	defer store.mx.Unlock()

	sessions, ok := store.symbols[symbol]
	if !ok || !isFreshSessions(sessions) {
		slog.Debug("get schedule sessions", "cache", "miss")

		var err error
		sessions, err := store.client.GetSchedule(store.ctx, symbol)
		if err != nil {
			return nil, err
		}
		sessions = normalizeSessions(sessions)
		store.symbols[symbol] = sessions
		return sessions, nil
	}
	slog.Debug("get schedule sessions", "cache", "hit")
	return sessions, nil
}

func normalizeSessions(sessions []types.Session) []types.Session {
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Start < sessions[j].Start && sessions[i].End < sessions[j].End
	})
	result := sessions[:0]
	last := func() int { return len(result) - 1 }
	for i, curr := range sessions {
		if i == 0 {
			result = append(result, curr)
			continue
		}
		if curr.Type == result[last()].Type {
			result[last()].End = curr.End
			continue
		}
		result = append(result, curr)
	}
	return result
}

func isFreshSessions(sessions []types.Session) bool {
	now := time.Now().Unix()
	for _, sess := range sessions {
		if sess.Start <= now && sess.End >= now {
			return true
		}
	}
	return false
}
