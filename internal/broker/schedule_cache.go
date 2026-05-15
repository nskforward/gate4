package broker

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/nskforward/gate4/pkg/types"
)

type ScheduleCache struct {
	clients map[Client]map[string][]types.Session
	mx      sync.Mutex
}

func NewScheduleCache() *ScheduleCache {
	return &ScheduleCache{
		clients: make(map[Client]map[string][]types.Session),
	}
}

func (cache *ScheduleCache) GetSessions(ctx context.Context, client Client, symbol string) ([]types.Session, error) {
	cache.mx.Lock()
	defer cache.mx.Unlock()

	symbols, ok := cache.clients[client]
	if !ok {
		symbols = make(map[string][]types.Session)
		cache.clients[client] = symbols
	}

	sessions, ok := symbols[symbol]
	if !ok || !isFreshSessions(sessions) {
		var err error
		sessions, err = client.GetSchedule(ctx, symbol)
		if err != nil {
			return nil, err
		}
		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].Start < sessions[j].Start && sessions[i].End < sessions[j].End
		})
		slog.Debug("get schedule sessions", "cache", "miss")
		symbols[symbol] = sessions
	} else {
		slog.Debug("get schedule sessions", "cache", "hit")
	}

	return sessions, nil
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
