package streams

import "sync/atomic"

type AtomicError struct {
	p atomic.Pointer[error]
}

func (a *AtomicError) Store(err error) {
	if err == nil {
		a.p.Store(nil)
		return
	}
	// Копируем, чтобы избежать алиасинга
	e := err
	a.p.Store(&e)
}
func (a *AtomicError) Load() error {
	ep := a.p.Load()
	if ep == nil {
		return nil
	}
	return *ep
}
