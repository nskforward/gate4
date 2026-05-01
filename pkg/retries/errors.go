package retries

import "errors"

var ErrMaxAttemptsExceeded = errors.New("max retry attempts exceeded")
var ErrStopped = errors.New("retry attemps stopped")
