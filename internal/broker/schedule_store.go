package broker

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/nskforward/gate4/internal/broker/types"
	"github.com/nskforward/gate4/pkg/finam"
)

type ScheduleStore struct {
	cache map[string][]types.Session
	mx    sync.Mutex
}

type FetchFunc func(ctx context.Context, account *Account, symbol string) ([]types.Session, error)

func NewScheduleStore() *ScheduleStore {
	return &ScheduleStore{
		cache: make(map[string][]types.Session),
	}
}

func (store *ScheduleStore) Sessions(ctx context.Context, client *finam.Client, symbol string) ([]types.Session, types.Session, error) {
	cache := "hit"
	var sessions []types.Session
	var ok bool
	var err error
	defer func() {
		slog.Debug("get sessions request", "broker", "finam", "symbol", symbol, "cache", cache, "sessions", len(sessions))
	}()

	sessions, ok = store.get(symbol)
	if ok {
		sess, ok := searchSession(sessions)
		if ok {
			return sessions, sess, nil
		}
	}

	cache = "miss"
	sessions, err = store.refresh(ctx, client, symbol)
	if err != nil {
		return nil, types.Session{}, fmt.Errorf("cannot fetch remote schedule sessions for symbol %s: %w", symbol, err)
	}

	sess, ok := searchSession(sessions)
	if !ok {
		return nil, types.Session{}, fmt.Errorf("cannot find the current session after refresh for symbol %s", symbol)
	}

	return sessions, sess, nil
}

func (store *ScheduleStore) refresh(ctx context.Context, client *finam.Client, symbol string) ([]types.Session, error) {
	result, err := client.GetSchedule(ctx, symbol)
	if err != nil {
		return nil, err
	}

	sessions := make([]types.Session, 0, len(result))
	for _, item := range result {
		sessType, err := sessionTypeCast(item.Type)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, types.Session{
			Type:  sessType,
			Start: item.Interval.StartTime.Seconds,
			End:   item.Interval.EndTime.Seconds,
		})
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Start < sessions[j].Start && sessions[i].End < sessions[j].End
	})
	store.save(symbol, sessions)
	return sessions, nil
}

func (store *ScheduleStore) get(key string) ([]types.Session, bool) {
	store.mx.Lock()
	defer store.mx.Unlock()
	items, ok := store.cache[key]
	return items, ok
}

func (store *ScheduleStore) save(key string, items []types.Session) {
	store.mx.Lock()
	defer store.mx.Unlock()
	store.cache[key] = items
}

func searchSession(sessions []types.Session) (types.Session, bool) {
	now := time.Now().Unix()

	for _, sess := range sessions {
		if now >= sess.Start && now <= sess.End {
			return sess, true
		}
	}
	return types.Session{}, false
}
