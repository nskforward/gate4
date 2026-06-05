package model

import (
	"errors"
	"time"
)

type Token struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

func (token Token) Validate() error {
	if token.ID == "" {
		return errors.New("token id cannot be empty")
	}
	if token.UserID == "" {
		return errors.New("token user_id cannot be empty")
	}
	if time.Since(token.Expires) > 0 {
		return errors.New("token expired")
	}
	return nil
}
