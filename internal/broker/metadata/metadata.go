package metadata

import "github.com/nskforward/gate4/pkg/types"

type Metadata struct {
	PositionStore *PositionStore
	QuoteStore    *QuoteStore
	ScheduleStore *ScheduleStore
}

func NewMetadata(client types.Client) (*Metadata, error) {
	positions, err := NewPositionStore(client)
	if err != nil {
		return nil, err
	}
	return &Metadata{
		PositionStore: positions,
		QuoteStore:    NewQuoteStore(client),
		ScheduleStore: NewScheduleStore(client),
	}, nil
}
