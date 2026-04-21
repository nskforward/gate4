package transport

import (
	"context"
	"fmt"

	"github.com/nskforward/gate4/pkg/pb"
)

func (s *Server) GetAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountResponse, error) {
	switch in.BrokerId {
	case "finam":
		client, err := s.finamAccounts.Get(in.AccountId)
		if err != nil {
			return nil, fmt.Errorf("finam client: %w", err)
		}
		resp, err := client.GetAccountInfo(ctx, in.AccountId)
		if err != nil {
			return nil, fmt.Errorf("finam communication error: %w", err)
		}
		return &pb.AccountResponse{
			BrokerId:  "finam",
			AccountId: resp.AccountId,
		}, nil
	}
	return nil, fmt.Errorf("broker '%s' not supported", in.BrokerId)
}
