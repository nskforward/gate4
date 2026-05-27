package di

import (
	"crypto/tls"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewTCPServer(apiServer *api.Server, tlsConfig *tls.Config) *grpc.Server {
	tcpSocket := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	pb.RegisterGate4Server(tcpSocket, apiServer)
	return tcpSocket
}
