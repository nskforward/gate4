package model

import (
	"errors"
	"time"
)

type Account struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Broker      string    `json:"broker"`
	Created     time.Time `json:"created"`
	Expires     time.Time `json:"expires"`
	Credentials []byte    `json:"credentials"`
}

func (account *Account) Validate() error {

	if account.ID == "" {
		return errors.New("account id cannot be empty")
	}

	if account.UserID == "" {
		return errors.New("user id cannot be empty")
	}

	if account.Broker == "" {
		return errors.New("broker name cannot be empty")
	}

	if len(account.Credentials) == 0 {
		return errors.New("credentials cannot be empty")
	}

	return nil
}
