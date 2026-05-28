package app

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/api/interceptor"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TCPServer struct {
	s    *grpc.Server
	l    net.Listener
	once sync.Once
}

func NewTCPServer(apiServer *api.Server, tlsConfig *tls.Config, tcpAddr string) *TCPServer {
	l, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		console.LogFatal("cannot listen tcp address", err)
	}
	s := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(
			interceptor.Logging,
			interceptor.Recovery,
		),
	)
	pb.RegisterGate4Server(s, apiServer)

	return &TCPServer{
		s: s,
		l: l,
	}
}

func (tcpServer *TCPServer) Serve() error {
	return tcpServer.s.Serve(tcpServer.l)
}

func (tcpServer *TCPServer) Close() {
	tcpServer.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stopped := make(chan struct{})
		go func() {
			tcpServer.s.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:

		case <-ctx.Done():
			tcpServer.s.Stop()
		}
	})
}
