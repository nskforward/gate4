package broker

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/nskforward/gate4/internal/broker/types"
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

func (store *ScheduleStore) CurrentSession(ctx context.Context, account *Account, symbol string, fetch FetchFunc) (types.Session, error) {
	refreshed := false
	sessions, ok := store.get(symbol)
	if !ok {
		newSessions, err := store.refresh(ctx, account, symbol, fetch)
		if err != nil {
			return types.Session{}, fmt.Errorf("cannot fetch remote schedule sessions for symbol %s: %w", symbol, err)
		}
		sessions = newSessions
		refreshed = true
	}

	sess, ok := searchSession(sessions)
	if !ok {
		if refreshed {
			return types.Session{}, fmt.Errorf("cannot find the current session for symbol %s", symbol)
		}
		newSessions, err := store.refresh(ctx, account, symbol, fetch)
		if err != nil {
			return types.Session{}, fmt.Errorf("cannot fetch remote schedule sessions for symbol %s: %w", symbol, err)
		}
		sessions = newSessions
		refreshed = true
	}

	sess, ok = searchSession(sessions)
	if !ok {
		return types.Session{}, fmt.Errorf("cannot find the current session for symbol %s", symbol)
	}

	return sess, nil
}

func (store *ScheduleStore) refresh(ctx context.Context, account *Account, symbol string, fetch FetchFunc) ([]types.Session, error) {
	newSessions, err := fetch(ctx, account, symbol)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch remote schedule sessions for symbol %s: %w", symbol, err)
	}
	sort.Slice(newSessions, func(i, j int) bool {
		return newSessions[i].Start < newSessions[j].Start && newSessions[i].End < newSessions[j].End
	})
	store.save(symbol, newSessions)
	return newSessions, nil
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
		if sess.Start >= now && sess.End <= now {
			return sess, true
		}
	}
	return types.Session{}, false
}
