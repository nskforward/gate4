include .env
export

.PHONY: build test

protoc:
	protoc --proto_path=proto --go_out=pkg/pb --go_opt=paths=source_relative --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative proto/*.proto

run_server:
	go run cmd/server/*

test:
	go test ./test/... -v

build_client:
	go build -o build/gate4 cmd/cli/*
