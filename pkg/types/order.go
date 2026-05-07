package types

type Order struct {
	Type   OrderType
	Symbol string
	Size   string // sign is direction (+ long / - short)
	Price  string // for Limit orders
}

func (order Order) IsShort() bool {
	return len(order.Size) > 0 && order.Size[0] == '-'
}

type OrderType string

const (
	OrderMarket OrderType = "MARKET"
	OrderLimit  OrderType = "LIMIT"
)

type OrderStatus string

const (
	OrderNew       OrderStatus = "new"
	OrderFilled    OrderStatus = "filled"
	OrderPartial   OrderStatus = "partial"
	OrderPending   OrderStatus = "pending"
	OrderCancelled OrderStatus = "cancelled"
	OrderRejected  OrderStatus = "rejected"
)
