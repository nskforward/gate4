package server

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

type UnixServer struct {
	socketPath string
	transport  *grpc.Server
}

/*
	interceptors := grpc.ChainUnaryInterceptor(
		interceptor.Logging,
		interceptor.Recovery,
	)
*/

func NewUnixServer(opts ...grpc.ServerOption) *UnixServer {
	s := &UnixServer{
		socketPath: filepath.Join(os.TempDir(), "gate4.sock"),
		transport:  grpc.NewServer(opts...),
	}
	return s
}

func (s *UnixServer) Register(handlers ...Registrable) {
	for _, h := range handlers {
		h.Register(s.transport)
	}
}

func (s *UnixServer) Start(ctx context.Context) error {
	l, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	defer os.Remove(s.socketPath)
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
	return nil
}
