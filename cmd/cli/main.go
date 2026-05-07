package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nskforward/gate4/internal/cli"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/transport/grpc/impl"
)

func main() {
	err := run(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "[error]:", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	adminClient, err := impl.NewAdminClient(cfg.Admin.ListenAddr)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	clientCLI := cli.NewClient(adminClient)
	return clientCLI.Run(ctx, os.Args)
}
