package finam

import (
	"fmt"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
)

type SessionType string

const (
	SessionClosed         SessionType = "CLOSED"          // nothing
	SessionOpeningAuction SessionType = "OPENING_AUCTION" // only limit
	SessionEarlyTrading   SessionType = "EARLY_TRADING"   // only limit
	SessionCoreTrading    SessionType = "CORE_TRADING"    // any
	SessionClosingAuction SessionType = "CLOSING_AUCTION" // only limit
	SessionLateTrading    SessionType = "LATE_TRADING"    // any
)

func NewSessionType(sessionType string) (SessionType, error) {
	switch sessionType {
	case "CLOSED":
		return SessionClosed, nil
	case "OPENING_AUCTION":
		return SessionOpeningAuction, nil
	case "EARLY_TRADING":
		return SessionEarlyTrading, nil
	case "CORE_TRADING":
		return SessionCoreTrading, nil
	case "CLOSING_AUCTION":
		return SessionClosingAuction, nil
	case "LATE_TRADING":
		return SessionLateTrading, nil
	default:
		return SessionType(""), fmt.Errorf("unknown sessionType: %s", sessionType)
	}
}

func (sess SessionType) Allow(orderType orders.OrderType) bool {
	switch sess {
	case SessionClosed:
		return false
	case SessionOpeningAuction:
		return orderType == orders.OrderType_ORDER_TYPE_LIMIT
	case SessionEarlyTrading:
		return orderType == orders.OrderType_ORDER_TYPE_LIMIT
	case SessionCoreTrading:
		return true
	case SessionClosingAuction:
		return orderType == orders.OrderType_ORDER_TYPE_LIMIT
	case SessionLateTrading:
		return true
	default:
		return false
	}
}
