package store

import "github.com/nskforward/gate4/pkg/pb"

type AccountProvider interface {
	Load() (map[string]*pb.Account, error)
	Save(map[string]*pb.Account) error
}
