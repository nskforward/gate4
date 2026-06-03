package model

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Blocked bool      `json:"blocked"`
	Created time.Time `json:"created"`
}

func (user User) Validate() error {
	if len(user.Email) < 3 || user.Email[0] == '@' || user.Email[len(user.Email)-1] == '@' || !strings.Contains(user.Email, "@") {
		return errors.New("invalid email")
	}

	if len(user.Name) < 2 {
		return errors.New("name length must be at least 2 characters")
	}

	return nil
}
