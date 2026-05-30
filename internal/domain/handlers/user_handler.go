package handlers

import (
	"github.com/nskforward/gate4/internal/domain/usecase"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	pb.UnimplementedUsersServer
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(serv *grpc.Server) {
	pb.RegisterUsersServer(serv, h)
}
