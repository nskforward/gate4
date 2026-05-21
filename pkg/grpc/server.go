package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	google "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	*google.Server
	addr     string
	OnListen func()
	OnStop   func()
}

func New(addr string, opt ...google.ServerOption) *Server {
	return &Server{
		addr:   addr,
		Server: google.NewServer(opt...),
	}
}

func NewWithTLS(addr string, tlsConfig *tls.Config) *Server {
	return &Server{
		addr:   addr,
		Server: google.NewServer(google.Creds(credentials.NewTLS(tlsConfig))),
	}
}

func (s *Server) Run(ctx context.Context) error {
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

func (s *Server) gracefulShutdown() {
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

func (s *Server) serve(listener net.Listener) chan error {
	errorc := make(chan error, 1)
	go func() {
		defer close(errorc)
		err := s.Server.Serve(listener)
		if err != nil && err != google.ErrServerStopped {
			errorc <- err
		}
	}()
	return errorc
}
