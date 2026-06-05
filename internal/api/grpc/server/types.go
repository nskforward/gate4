package server

import "google.golang.org/grpc"

type Registrable interface {
	Register(s *grpc.Server)
}
