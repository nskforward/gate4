package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGatewayServer
	addr      string
	transport *grpc.Server
}

func NewServer(addr string) *Server {
	transportServer := grpc.NewServer()
	s := &Server{
		addr:      addr,
		transport: transportServer,
	}
	pb.RegisterGatewayServer(transportServer, s)
	return s
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net.Listen error: %w", err)
	}

	slog.Info("start server", "addr", s.addr)

	errorc := s.serve(listener)

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errorc:
		return err

	case <-sigCtx.Done():
		s.gracefulShutdown()
		return nil
	}

}

func (s *Server) gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
}

func (s *Server) serve(listener net.Listener) chan error {
	errorc := make(chan error, 1)
	go func() {
		defer close(errorc)
		err := s.transport.Serve(listener)
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()
	return errorc
}
