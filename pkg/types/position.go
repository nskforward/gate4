package types

type Position struct {
	Symbol string
	Price  string
	Size   string // sign is direction
}

func (pos Position) IsEmpty() bool {
	return pos.Size == ""
}
