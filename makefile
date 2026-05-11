include .env
export

.PHONY: build test

protoc:
	protoc --proto_path=proto --go_out=pkg/pb --go_opt=paths=source_relative --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative proto/*.proto

run:
	go run cmd/server/* \
	 --store-dir=data \
	 --ssl-ca-key=data/ssl/ca.key \
	 --ssl-ca-cert=data/ssl/ca.crt \
	 --ssl-server-key=data/ssl/server.key \
	 --ssl-server-cert=data/ssl/server.crt \

test:
	go test ./test/... -v

client:
	go build -o build/gate4 cmd/cli/*
