package finam

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func connect(addr string) (conn *grpc.ClientConn, err error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})),
		grpc.WithIdleTimeout(5*time.Minute),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
	)
}
