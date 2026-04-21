package transport

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/pkg/pb"
)

func (s *Server) GetAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	return nil, errors.New("not implemented")
}
