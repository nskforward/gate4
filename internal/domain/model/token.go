package model

import "time"

type Token struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	Role    Role      `json:"role"`
	Expires time.Time `json:"expires"`
}

func (token Token) Validate() error {

}
