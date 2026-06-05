package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type TCPServer struct {
	addr      string
	transport *grpc.Server
}

/*
	grpc.Creds(credentials.NewTLS(tlsConfig)),
	interceptors := grpc.ChainUnaryInterceptor(
		interceptor.Logging,
		interceptor.Recovery,
	)
*/

func NewTCPServer(addr string, opts ...grpc.ServerOption) *TCPServer {
	return &TCPServer{
		addr:      addr,
		transport: grpc.NewServer(opts...),
	}
}

func (s *TCPServer) Register(handlers ...Registrable) {
	for _, h := range handlers {
		h.Register(s.transport)
	}
}

func (s *TCPServer) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	err = s.transport.Serve(l)
	if err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

func (s *TCPServer) Stop(ctx context.Context) error {
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
