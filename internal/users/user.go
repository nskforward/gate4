package users

import (
	"time"

	"github.com/nskforward/gate4/internal/brokers"
)

type User struct {
	ID         string           `json:"id"`
	BrokerID   brokers.BrokerID `json:"broker_id"`
	AccountID  string           `json:"account_id"`
	Secret     string           `json:"secret"`
	ValidUntil time.Time        `json:"valid_until"`
	Blocked    bool             `json:"blocked"`
}
