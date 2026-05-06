package types

type Order struct {
	Type      OrderType
	Timestamp int64
	Symbol    string
	Size      string // sign is direction (+ long / - short)
	Price     string // for Limit orders
}

func (order Order) IsShort() bool {
	return len(order.Size) > 0 && order.Size[0] == '-'
}

type OrderType string

const (
	OrderMarket OrderType = "MARKET"
	OrderLimit  OrderType = "LIMIT"
)
