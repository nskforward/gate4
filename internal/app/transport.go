package app

import (
	"crypto/tls"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetStartUnixServer(apiServer *api.Server) func() error {
	socketPath := filepath.Join(os.TempDir(), "gate4.sock")

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		slog.Error("cannot listen unix socket address", "reason", err.Error())
		os.Exit(1)
	}

	transport := grpc.NewServer()
	pb.RegisterGate4Server(transport, apiServer)

	return func() error {
		defer os.Remove(socketPath)
		return transport.Serve(listener)
	}
}

func GetStartTCPServer(apiServer *api.Server, tlsConfig *tls.Config, tcpAddr string) func() error {
	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		slog.Error("cannot listen tcp address", "reason", err)
		os.Exit(1)
	}

	transport := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	pb.RegisterGate4Server(transport, apiServer)

	return func() error {
		return transport.Serve(listener)
	}
}
