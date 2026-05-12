package retries

/*
func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	cfg := Config{
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		MaxAttempts:   3,
		JitterFactor:  0,
	}

	expected := "success"
	retry := NewRetry[string](cfg, func() (string, error) {
		return expected, nil
	})

	ctx := context.Background()
	result, err := retry.Do(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRetry_SuccessAfterRetries(t *testing.T) {
	cfg := Config{
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      100 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   5,
		JitterFactor:  0,
	}

	var attempts atomic.Int32
	expected := "success"
	retry := NewRetry[string](cfg, func() (string, error) {
		attempt := attempts.Add(1)
		if attempt < 3 {
			return "", errors.New("temporary error")
		}
		return expected, nil
	})

	ctx := context.Background()
	result, err := retry.Do(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestRetry_MaxAttemptsExceeded(t *testing.T) {
	cfg := Config{
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      100 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   3,
		JitterFactor:  0,
	}

	expectedErr := errors.New("permanent error")
	retry := NewRetry[string](cfg, func() (string, error) {
		return "", expectedErr
	})

	ctx := context.Background()
	_, err := retry.Do(ctx)

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestRetry_ContextCanceled(t *testing.T) {
	cfg := Config{
		InitialDelay:  1 * time.Second,
		MaxDelay:      2 * time.Second,
		BackoffFactor: 2.0,
		MaxAttempts:   0, // infinite
		JitterFactor:  0,
	}

	retry := NewRetry[string](cfg, func() (string, error) {
		return "", errors.New("error")
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := retry.Do(ctx)

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRetry_ExponentialBackoff(t *testing.T) {
	cfg := Config{
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      100 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   0,
		JitterFactor:  0,
	}

	var attempts atomic.Int32

	retry := NewRetry[string](cfg, func() (string, error) {
		attempt := attempts.Add(1)
		if attempt < 3 {
			return "", errors.New("error")
		}
		return "", errors.New("stop")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()
	retry.Do(ctx)
	elapsed := time.Since(start)

	// With initial=10ms, factor=2, attempts: 10ms, 20ms, 40ms = ~70ms
	// Allow some margin
	if elapsed < 60*time.Millisecond {
		t.Errorf("expected at least 60ms delay, got %v", elapsed)
	}
}

func TestRetry_MaxDelayCapped(t *testing.T) {
	cfg := Config{
		InitialDelay:  1 * time.Millisecond,
		MaxDelay:      20 * time.Millisecond,
		BackoffFactor: 10.0, // very aggressive
		MaxAttempts:   5,
		JitterFactor:  0,
	}

	var attempts atomic.Int32

	retry := NewRetry[string](cfg, func() (string, error) {
		attempts.Add(1)
		return "", errors.New("error")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	retry.Do(ctx)
	elapsed := time.Since(start)

	// Should be capped at MaxDelay (20ms) per attempt, 5 attempts = ~100ms max
	// But actual delay is per attempt, so with 5 attempts at 20ms max = ~100ms
	if elapsed > 150*time.Millisecond {
		t.Errorf("expected delay around 100ms, got %v", elapsed)
	}
	if attempts.Load() != 5 {
		t.Errorf("expected 5 attempts, got %d", attempts.Load())
	}
}

func TestRetry_JitterFactor(t *testing.T) {
	cfg := Config{
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      500 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   0,
		JitterFactor:  0.5, // 50% jitter
	}

	var attempts atomic.Int32

	retry := NewRetry[string](cfg, func() (string, error) {
		if attempts.Add(1) < 2 {
			return "", errors.New("error")
		}
		return "done", nil // успех на 3й попытке
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	retry.Do(ctx)
	elapsed := time.Since(start)

	// With jitter, delay for first retry should be in range [100ms, 150ms]
	if elapsed < 95*time.Millisecond || elapsed > 160*time.Millisecond {
		t.Errorf("expected delay between 95-160ms, got %v", elapsed)
	}
}

func TestRetry_JitterNotExceedMaxDelay(t *testing.T) {
	cfg := Config{
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      150 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   5,
		JitterFactor:  1.0, // 100% jitter could push over max
	}

	var attempts atomic.Int32

	retry := NewRetry[string](cfg, func() (string, error) {
		attempts.Add(1)
		return "", errors.New("error")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()
	retry.Do(ctx)
	elapsed := time.Since(start)

	// Should not exceed MaxDelay + MaxDelay*JitterFactor = 300ms per attempt
	// With 5 attempts should be around 500-600ms
	if elapsed > 800*time.Millisecond {
		t.Errorf("expected delay under 800ms, got %v", elapsed)
	}
}

func TestRetry_InfiniteAttempts(t *testing.T) {
	cfg := Config{
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      20 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   0, // infinite
		JitterFactor:  0,
	}

	var attempts atomic.Int32

	retry := NewRetry[string](cfg, func() (string, error) {
		attempts.Add(1)
		return "", errors.New("error")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	retry.Do(ctx)

	// Should have many attempts before timeout
	if attempts.Load() < 3 {
		t.Errorf("expected at least 3 attempts, got %d", attempts.Load())
	}
}

func TestRetry_ContextCanceledDuringDelay(t *testing.T) {
	cfg := Config{
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      500 * time.Millisecond,
		BackoffFactor: 2.0,
		MaxAttempts:   0,
		JitterFactor:  0,
	}

	var attempts atomic.Int32

	retry := NewRetry[string](cfg, func() (string, error) {
		attempts.Add(1)
		return "", errors.New("error")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := retry.Do(ctx)
	elapsed := time.Since(start)

	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
	// Should have made only 1 attempt and stopped
	if attempts.Load() > 2 {
		t.Errorf("expected 1-2 attempts, got %d", attempts.Load())
	}
	if elapsed > 100*time.Millisecond {
		t.Errorf("expected quick stop, got %v", elapsed)
	}
}
*/
