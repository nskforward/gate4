package retries

import (
	"context"
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
		MaxDelay:     2 * time.Second,
		Backoff:      1.5,
		MaxJitter:    time.Second,
	}

	t1 := time.Now()
	for attempt := range r.Range(ctx) {
		fmt.Println(time.Since(t1), "attempt", attempt)
		t1 = time.Now()
	}

	fmt.Println(time.Since(t1), "stopped the test with error:", r.Err())
}
