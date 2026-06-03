package model

import "errors"

type Role string

const (
	Screener Role = "screener"
	Trader   Role = "trader"
	Admin    Role = "admin"
)

func (role Role) String() string {
	switch role {
	case Screener:
		return "screener"
	case Trader:
		return "trader"
	case Admin:
		return "admin"
	default:
		return ""
	}
}

func (role Role) Validate() error {
	if role.String() == "" {
		return errors.New("unknown role")
	}
	return nil
}
