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

func TestTopic1(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	toptic := newTopic("1", func(ctx context.Context, f func(data int) bool) error {
		defer fmt.Println("publisher stopped")

		time.Sleep(time.Second)
		fmt.Println("start publisher")

		for i := range 10 {
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

	go func() {
		defer wg.Done()
		for v := range s1.Range() {
			fmt.Println("s1 <- ", v)
		}
		fmt.Println("s1 stopped to read")
		s1.Close()
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
}
