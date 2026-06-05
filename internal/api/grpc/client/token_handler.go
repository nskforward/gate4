package client

import (
	"context"
	"time"

	"github.com/nskforward/gate4/internal/api/grpc/common"
	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type TokenHandler struct {
	client pb.TokensClient
}

func NewTokenHandler(conn *grpc.ClientConn) *TokenHandler {
	return &TokenHandler{
		client: pb.NewTokensClient(conn),
	}
}

func (h *TokenHandler) ListUserTokens(ctx context.Context, userID string) ([]model.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	result, err := h.client.ListUserTokens(ctx, &pb.UserID{UserId: userID})
	if err != nil {
		return nil, wrapError(err)
	}
	return common.ConvertOutTokens(result.Tokens), nil
}

func (h *TokenHandler) CreateToken(ctx context.Context, token *model.Token) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	result, err := h.client.CreateToken(ctx, common.ConvertInToken(*token))
	if err != nil {
		return wrapError(err)
	}

	token.ID = result.Id
	token.Created = time.Unix(result.Created, 0)

	return nil
}

func (h *TokenHandler) DeleteToken(ctx context.Context, tokenID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := h.client.DeleteToken(ctx, &pb.TokenID{TokenId: tokenID})
	return wrapError(err)
}
