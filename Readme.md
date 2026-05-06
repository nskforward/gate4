# GATE4

## Architecture

proto impl -> broker -> finam client

## TODO

- pkg/finam must returns broker/types instead of native finam
- broker/broker must return broker/types instead of pkg/pb
- pkg/finam client must implement broker/client interface
- transport/grpc must cast broker/types to pkg/pb