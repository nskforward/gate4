package broker

import (
	"context"
	"sync"

	"github.com/nskforward/gate4/pkg/pubsub"
	"github.com/nskforward/gate4/pkg/types"
)

type quoteStore struct {
	ctx     context.Context
	cancel  context.CancelFunc
	client  Client
	symbols map[string]*pubsub.Topic[types.Quote]
	mx      sync.Mutex
}

func newQuoteStore(client Client) *quoteStore {
	ctx, cancel := context.WithCancel(context.Background())
	return &quoteStore{
		ctx:     ctx,
		cancel:  cancel,
		client:  client,
		symbols: make(map[string]*pubsub.Topic[types.Quote]),
	}
}

func (store *quoteStore) GetLast(symbol string) (types.Quote, error) {
	store.mx.Lock()
	topic, ok := store.symbols[symbol]
	store.mx.Unlock()
	if ok {
		q := topic.Last()
		if q != nil {
			return *q, nil
		}
	}
	return store.client.GetLastQuote(store.ctx, symbol)
}

func (store *quoteStore) Subscribe(ctx context.Context, symbol string, f func(types.Quote) error) error {
	store.mx.Lock()
	defer store.mx.Unlock()

	topic, ok := store.symbols[symbol]
	if !ok {
		topic = pubsub.NewTopic[types.Quote](1)
		store.symbols[symbol] = topic
		// TODO start producer !!!
	}
	sub, err := topic.Subscribe(ctx)
	if err != nil {
		return err
	}
	for q := range sub.Range() {
		err := f(q)
		if err != nil {
			return err
		}
	}
	return nil
}
