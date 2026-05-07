package types

import "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"

type Session struct {
	Type  SessionType
	Start int64 // unix timestamp in seconds
	End   int64 // unix timestamp in seconds
}

type SessionType string

const (
	SessionUnspecified SessionType = "UNSPECIFIED"
	SessionClosed      SessionType = "CLOSED"
	SessionMain        SessionType = "MAIN"
	SessionPremarket   SessionType = "PREMARKET"
	SessionPostmarket  SessionType = "POSTMARKET"
)

func (sess SessionType) Allow(orderType orders.OrderType) bool {
	switch sess {
	case SessionMain:
		return true
	case SessionPremarket, SessionPostmarket:
		return orderType == orders.OrderType_ORDER_TYPE_LIMIT
	default:
		return false
	}
}
