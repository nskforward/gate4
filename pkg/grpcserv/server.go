package grpcserv

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	addr     string
	OnListen func()
	OnStop   func()
}

func New(addr string) *GRPCServer {
	return &GRPCServer{
		addr:   addr,
		Server: grpc.NewServer(),
	}
}

func (s *GRPCServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net.Listen error: %w", err)
	}

	if s.OnListen != nil {
		s.OnListen()
	}

	errorc := s.serve(listener)

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errorc:
		if s.OnStop != nil {
			s.OnStop()
		}
		return err

	case <-sigCtx.Done():
		s.gracefulShutdown()
		if s.OnStop != nil {
			s.OnStop()
		}
		return nil
	}
}

func (s *GRPCServer) gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		s.Server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-ctx.Done():
		s.Server.Stop()
	}
}

func (s *GRPCServer) serve(listener net.Listener) chan error {
	errorc := make(chan error, 1)
	go func() {
		defer close(errorc)
		err := s.Server.Serve(listener)
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()
	return errorc
}
