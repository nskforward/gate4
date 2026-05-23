package brokers

import "fmt"

type BrokerID string

const (
	FINAM BrokerID = "finam"
)

func (id BrokerID) Validate() error {
	switch id {
	case FINAM:
		return nil
	default:
		return fmt.Errorf("unknown broker id")
	}
}
