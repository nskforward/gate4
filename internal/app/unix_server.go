package app

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
)

type UnixServer struct {
	s          *grpc.Server
	l          net.Listener
	socketPath string
	once       sync.Once
}

func NewUnixServer(apiServer *api.Server) *UnixServer {
	socketPath := filepath.Join(os.TempDir(), "gate4.sock")
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		console.LogFatal("cannot listen unix socket address", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		api.LoggingInterceptor,
		api.RecoveryInterceptor,
	))

	pb.RegisterGate4Server(s, apiServer)

	return &UnixServer{
		s:          s,
		l:          l,
		socketPath: socketPath,
	}
}

func (unixServer *UnixServer) Serve() error {
	defer os.Remove(unixServer.socketPath)
	return unixServer.s.Serve(unixServer.l)
}

func (unixServer *UnixServer) Close() {
	unixServer.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stopped := make(chan struct{})
		go func() {
			unixServer.s.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:

		case <-ctx.Done():
			unixServer.s.Stop()
		}
	})
}
