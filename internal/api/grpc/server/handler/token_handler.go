package handler

import (
	"context"

	"github.com/nskforward/gate4/internal/api/grpc/common"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type TokenHandler struct {
	pb.UnimplementedTokensServer
	tokenService *service.TokenService
}

func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
	}
}

func (h *TokenHandler) Register(s *grpc.Server) {
	pb.RegisterTokensServer(s, h)
}

func (h *TokenHandler) ListUserTokens(ctx context.Context, req *pb.UserID) (*pb.TokenList, error) {
	tokens, err := h.tokenService.ListUserTokens(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.TokenList{
		Tokens: common.ConvertInTokens(tokens),
	}, nil
}

func (h *TokenHandler) CreateToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {
	token := common.ConvertOutToken(req)
	err := h.tokenService.CreateToken(ctx, &token)
	if err != nil {
		return nil, err
	}
	return common.ConvertInToken(token), nil
}

func (h *TokenHandler) DeleteToken(ctx context.Context, req *pb.TokenID) (*pb.EmptyMessage, error) {
	return &pb.EmptyMessage{}, h.tokenService.DeleteToken(ctx, req.TokenId)
}

func (h *TokenHandler) Whoami(ctx context.Context, req *pb.TokenID) (*pb.User, error) {
	user, err := h.tokenService.FindUser(ctx, req.TokenId)
	if err != nil {
		return nil, err
	}
	return common.ConvertInUser(user), nil
}
