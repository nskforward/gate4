package stream

import "sync/atomic"

type AtomicError struct {
	v atomic.Value
}

func (e *AtomicError) Store(err error) {
	if err == nil {
		e.v.Store(nil)
		return
	}
	e.v.Store(struct{ err error }{err}) // Wrap non-nil errors
}

func (e *AtomicError) Load() error {
	v := e.v.Load()
	if v == nil {
		return nil
	}
	return v.(struct{ err error }).err
}
