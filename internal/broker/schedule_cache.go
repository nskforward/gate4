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
		slog.Debug("get schedule sessions", "cache", "miss")

		var err error
		sessions, err = client.GetSchedule(ctx, symbol)
		if err != nil {
			return nil, err
		}
		sessions = cache.normalize(sessions)
		symbols[symbol] = sessions
		return sessions, nil
	}

	slog.Debug("get schedule sessions", "cache", "hit")
	return sessions, nil
}

func (cache *ScheduleCache) normalize(sessions []types.Session) []types.Session {
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
