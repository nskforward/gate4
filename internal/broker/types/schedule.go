package types

import "log/slog"

type Session struct {
	Type  SessionType
	Start int64 // unix timestamp in seconds
	End   int64 // unix timestamp in seconds
}

type SessionType string

const (
	SessionClosed     SessionType = "CLOSED"
	SessionMain       SessionType = "MAIN"
	SessionPremarket  SessionType = "PREMARKET"
	SessionPostmarket SessionType = "POSTMARKET"
)

func (sessType SessionType) Allow(orderType OrderType) bool {
	switch sessType {
	case SessionClosed:
		return false
	case SessionMain:
		return true
	case SessionPremarket, SessionPostmarket:
		return orderType == OrderLimit
	default:
		slog.Error("unknown order type", "have", orderType)
		return false
	}
}
