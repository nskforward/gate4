package tools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsGRPCCancelled(err error) bool {
	st, ok := status.FromError(err)
	if ok {
		if st.Code() == codes.Canceled {
			return true
		}
	}
	return false
}
