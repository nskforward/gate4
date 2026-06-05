package client

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.Unavailable {
		return errors.New("server unavailable")
	}
	return err
}
