package store

import (
	"sync"

	"github.com/nskforward/gate4/pkg/pb"
)

type BrokerStore struct {
	items map[string]*pb.Broker
	mx    sync.Mutex
}

func NewBrokerStore() *BrokerStore {
	return &BrokerStore{
		items: make(map[string]*pb.Broker),
	}
}

func (s *BrokerStore) List() []*pb.Broker {
	s.mx.Lock()
	defer s.mx.Unlock()
	result := make([]*pb.Broker, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, &pb.Broker{
			Id:      item.Id,
			Address: item.Address,
		})
	}
	return result
}
