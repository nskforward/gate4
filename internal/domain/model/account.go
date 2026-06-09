package model

import "time"

type Account struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	Broker  string    `json:"broker"`
	Created time.Time `json:"created"`
	Creds   []byte    `json:"creds"`
}
