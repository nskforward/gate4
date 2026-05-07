package types

type Quote struct {
	Timestamp int64
	Symbol    string
	Ask       QuoteLevel
	Bid       QuoteLevel
}

type QuoteLevel struct {
	Price string
	Size  string
}
