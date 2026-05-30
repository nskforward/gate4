package handlers

import (
	"github.com/nskforward/gate4/internal/domain/usecases"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UserHandler struct {
	pb.UnimplementedUsersServer
	userUsecases *usecases.UserUsecases
}

func NewUserHandler(userUsecases *usecases.UserUsecases) *UserHandler {
	return &UserHandler{
		userUsecases: userUsecases,
	}
}

func (h *UserHandler) Register(serv *grpc.Server) {
	pb.RegisterUsersServer(serv, h)
}
