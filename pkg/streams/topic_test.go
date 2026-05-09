package streams

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"
)

func TestTopicSimple(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	toptic := newTopic(initOpts(nil), "test", func(ctx context.Context, publish func(data int) bool) error {
		defer fmt.Println("publisher stopped")

		time.Sleep(10 * time.Millisecond)

		fmt.Println("start publisher")

		for i := range 10 {
			if !publish(i) {
				fmt.Println("stop publisher because no subscribers")
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(2)

	s1 := toptic.Subscribe(context.Background())
	s2 := toptic.Subscribe(context.Background())
	defer s1.Close()
	defer s2.Close()

	go func() {
		defer wg.Done()
		for v := range s1.Range() {
			fmt.Println("s1 <- ", v)
		}
		checkExitStatus("s1", s1)
	}()

	go func() {
		defer wg.Done()
		for v := range s2.Range() {
			fmt.Println("s2 <- ", v)
		}
		checkExitStatus("s2", s2)
	}()

	wg.Wait()
}

func TestTopicRetry(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	atttempt := 0

	toptic := newTopic(initOpts(nil), "test", func(ctx context.Context, publish func(data int) bool) error {
		defer fmt.Println("publisher stopped")

		time.Sleep(100 * time.Millisecond)
		fmt.Println("start publisher")

		for i := range 5 {
			if i == 0 && atttempt < 3 {
				atttempt++
				return fmt.Errorf("publisher test disconnection")
			}
			if !publish(i) {
				fmt.Println("stop publisher because no subscribers")
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(2)

	s1 := toptic.Subscribe(context.Background())
	s2 := toptic.Subscribe(context.Background())
	defer s1.Close()
	defer s2.Close()

	go func() {
		defer wg.Done()
		for v := range s1.Range() {
			fmt.Println("s1 <- ", v)
		}
		checkExitStatus("s1", s1)
	}()

	go func() {
		defer wg.Done()
		for v := range s2.Range() {
			fmt.Println("s2 <- ", v)
		}
		checkExitStatus("s2", s2)
	}()

	wg.Wait()
}

func TestTopicOneStreamCancel(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	toptic := newTopic(initOpts(nil), "test", func(ctx context.Context, publish func(data int) bool) error {
		defer fmt.Println("publisher stopped")

		time.Sleep(100 * time.Millisecond)
		fmt.Println("start publisher")

		for i := range 5 {
			if !publish(i) {
				fmt.Println("stop publisher because no subscribers")
				break
			}
			time.Sleep(400 * time.Millisecond)
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s1 := toptic.Subscribe(ctx)
	s2 := toptic.Subscribe(context.Background())
	defer s1.Close()
	defer s2.Close()

	go func() {
		defer wg.Done()
		for v := range s1.Range() {
			fmt.Println("s1 <- ", v)
		}
		checkExitStatus("s1", s1)
	}()

	go func() {
		defer wg.Done()
		for v := range s2.Range() {
			fmt.Println("s2 <- ", v)
		}
		checkExitStatus("s2", s2)
	}()

	wg.Wait()
}

func checkExitStatus(name string, s *Stream[int]) {
	err := s.Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println(name, "context deadline")
		} else {
			fmt.Println(name, "error:", err)
		}
	} else {
		fmt.Println(name, "normal close")
	}
}
