package model

import "time"

type User struct {
	ID      string
	Email   string
	Blocked bool
	Created time.Time
}
