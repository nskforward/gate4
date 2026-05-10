package streams

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	s := newStream(context.Background(), 1, func(s *Stream[int]) {
		fmt.Println("[info] stream destroyed")
	})
	defer s.Close()

	fmt.Println("[info] stream created")

	go generator(s)

	for v := range s.Range() {
		fmt.Println("[info] read:", v)
	}

	if err := s.Err(); err != nil {
		fmt.Println("stream has error:", err)
	}

	fmt.Println("finish")
}

func TestStreamTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s := newStream(ctx, 1, func(s *Stream[int]) {
		fmt.Println("[info] stream destroyed")
	})
	defer s.Close()

	fmt.Println("[info] stream created")

	go generator(s)

	for v := range s.Range() {
		fmt.Println("[info] read:", v)
	}

	if err := s.Err(); err != nil {
		fmt.Println("stream has error:", err)
	}

	fmt.Println("finish")
}

func generator(s *Stream[int]) {
	defer fmt.Println("generator stopped")
	fmt.Println("start generator")
	for i := range 10 {
		s.notify(i)
		select {
		case <-s.ctx.Done():
			s.err.Store(s.ctx.Err())
			return
		case <-time.After(200 * time.Millisecond):
			continue
		}
	}
	s.Close()
}
