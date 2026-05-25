package retry

type Attempt struct {
	r   *Retry
	cfg *config
}

func (a *Attempt) Num() int {
	return a.r.currentAttempt
}

func (a *Attempt) Success() {
	a.r.currentAttempt = 0
	a.r.currentDelay = a.cfg.InitialDelay
	a.r.err = nil
}

func (a *Attempt) Fail(err error) {
	a.r.err = err
}
