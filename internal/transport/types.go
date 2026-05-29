package transport

import "google.golang.org/grpc"

type ServiceRegistrar interface {
	Register(server *grpc.Server)
}
