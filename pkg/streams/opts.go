package streams

import (
	"math/rand"
	"time"
)

type Opts struct {
	BufferSize       int
	MaxRetryAttempts int32
	RetryWait        func(attempt int32) time.Duration
}

var (
	defaultOpt = Opts{
		BufferSize:       1,
		MaxRetryAttempts: 5,
		RetryWait: func(attempt int32) time.Duration {
			jitter := int32(rand.Intn(1000))
			mills := attempt*1000 + jitter
			return time.Duration(mills) + time.Millisecond
		},
	}
)

func initOpts(in *Opts) *Opts {
	if in == nil {
		return &defaultOpt
	}
	if in.BufferSize < 1 {
		in.BufferSize = defaultOpt.BufferSize
	}
	if in.MaxRetryAttempts < 0 {
		in.MaxRetryAttempts = 0
	}
	if in.MaxRetryAttempts > 0 && in.RetryWait == nil {
		in.RetryWait = defaultOpt.RetryWait
	}
	return in
}
