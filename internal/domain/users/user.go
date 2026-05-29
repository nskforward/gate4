package users

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	BrokerID  BrokerID  `json:"broker_id"`
	AccountID string    `json:"account_id"`
	Secret    string    `json:"secret"`
	Created   time.Time `json:"created"`
	Expires   time.Time `json:"expires"`
	Blocked   bool      `json:"blocked"`
}
