package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/nskforward/gate4/internal/domain/service"
)

type HTTPServer struct {
	transport    *http.Server
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	return errors.ErrUnsupported
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return errors.ErrUnsupported
}
