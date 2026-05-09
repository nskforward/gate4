package streams

import (
	"context"
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

	toptic := newTopic("test", func(ctx context.Context, f func(data int) bool) error {
		defer fmt.Println("publisher stopped")

		time.Sleep(10 * time.Millisecond)
		fmt.Println("start publisher")

		for i := range 10 {
			if !f(i) {
				fmt.Println("stop publisher because no subscribers")
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(2)

	s1 := toptic.Subscribe()
	s2 := toptic.Subscribe()
	defer s1.Close()
	defer s2.Close()

	go func() {
		defer wg.Done()
		for v := range s1.Range() {
			fmt.Println("s1 <- ", v)
		}
		fmt.Println("s1 stopped to read")
	}()

	go func() {
		defer wg.Done()
		for v := range s2.Range() {
			fmt.Println("s2 <- ", v)
		}
		fmt.Println("s2 stopped to read")
		s2.Close()
	}()

	wg.Wait()

	if err := s1.Err(); err != nil {
		fmt.Println("s1 error:", err)
	}
	if err := s2.Err(); err != nil {
		fmt.Println("s2 error:", err)
	}
}

func TestTopicRetry(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	atttempt := 0

	toptic := newTopic("test", func(ctx context.Context, f func(data int) bool) error {
		defer fmt.Println("publisher stopped")

		time.Sleep(100 * time.Millisecond)
		fmt.Println("start publisher")

		for i := range 5 {
			if i == 3 && atttempt < 5 {
				atttempt++
				return fmt.Errorf("publisher test disconnection")
			}
			if atttempt > 0 && i == 0 && atttempt < 3 {
				atttempt++
				return fmt.Errorf("publisher test disconnection")
			}
			if !f(i) {
				fmt.Println("stop publisher because no subscribers")
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(2)

	s1 := toptic.Subscribe()
	s2 := toptic.Subscribe()
	defer s1.Close()
	defer s2.Close()

	go func() {
		defer wg.Done()
		for v := range s1.Range() {
			fmt.Println("s1 <- ", v)
		}
		fmt.Println("s1 stopped to read")
	}()

	go func() {
		defer wg.Done()
		for v := range s2.Range() {
			fmt.Println("s2 <- ", v)
		}
		fmt.Println("s2 stopped to read")
	}()

	wg.Wait()

	if err := s1.Err(); err != nil {
		fmt.Println("s1 error:", err)
	}
	if err := s2.Err(); err != nil {
		fmt.Println("s2 error:", err)
	}
}
