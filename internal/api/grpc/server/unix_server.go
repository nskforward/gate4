package server

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/nskforward/gate4/internal/api/grpc/server/handler"
	"github.com/nskforward/gate4/internal/api/grpc/server/interceptor"
	"google.golang.org/grpc"
)

type UnixServer struct {
	socketPath string
	transport  *grpc.Server
}

func NewUnixServer(userHandler *handler.UserHandler, tokenHandler *handler.TokenHandler) *UnixServer {
	s := &UnixServer{
		socketPath: filepath.Join(os.TempDir(), "gate4.sock"),
		transport: grpc.NewServer(grpc.ChainUnaryInterceptor(
			interceptor.Logging,
			interceptor.Recovery,
		)),
	}

	userHandler.Register(s.transport)
	tokenHandler.Register(s.transport)

	return s
}

func (s *UnixServer) Start(ctx context.Context) error {
	l, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	err = s.transport.Serve(l)
	if err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

func (s *UnixServer) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.transport.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:

	case <-ctx.Done():
		s.transport.Stop()
	}

	os.Remove(s.socketPath)

	return nil
}
