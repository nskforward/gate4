package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := Retry{
		MaxAttempts:  5,
		InitialDelay: time.Second,
		MaxDelay:     10 * time.Second,
		Backoff:      1.5,
		MaxJitter:    time.Second,
	}

	t1 := time.Now()
	for attempt := range r.Range(ctx) {
		fmt.Println(time.Since(t1), "attempt", attempt.Num())
		t1 = time.Now()
		if attempt.Num() == 3 {
			attempt.Success()
		} else {
			attempt.Fail(errors.New("the test error"))
		}
	}

	fmt.Println(time.Since(t1), "stopped the test with error:", r.Err())
}
