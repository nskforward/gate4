package types

type Quote struct {
	Symbol    string
	Timestamp int64
	AskPrice  []string
	AskSize   []string
	BidPrice  []string
	BidSize   []string
}
