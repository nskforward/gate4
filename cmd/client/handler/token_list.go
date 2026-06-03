package handler

import (
	"context"
	"errors"

	"github.com/nskforward/gate4/internal/api/grpc/client"
	"github.com/nskforward/gate4/pkg/console/router"
)

func ListTokens(c *client.Client) router.Handler {
	return func(ctx context.Context, args []string) error {

		if len(args) < 1 {
			return errors.New("requires 1 argument")
		}

		var userID string

		if len(args) == 1 {
			userID = args[0]
			args = args[1:]
		}

		_ = userID

		return errors.ErrUnsupported
	}
}
