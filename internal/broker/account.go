package broker

import (
	"fmt"
	"time"

	"github.com/nskforward/gate4/pkg/pb"
)

type Account struct {
	ID         string
	Broker     string
	Secret     string
	ValidUntil time.Time
}

func (account *Account) Key() string {
	return fmt.Sprintf("%s.%s", account.Broker, account.ID)
}

func ImportAccount(in *pb.Account) *Account {
	return &Account{
		ID:         in.Id,
		Broker:     in.BrokerId,
		Secret:     in.Secret,
		ValidUntil: time.Unix(in.ValidUntil, 0),
	}
}

func ExportAccounts(list []*Account) []*pb.Account {
	result := make([]*pb.Account, 0, len(list))
	for _, account := range list {
		result = append(result, &pb.Account{
			BrokerId:   account.Broker,
			Id:         account.ID,
			ValidUntil: account.ValidUntil.Unix(),
		})
	}
	return result
}
