package streams

import (
	"fmt"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	s := newStream(func(s *Stream[int]) {
		fmt.Println("[info] stream destroyed")
	})
	fmt.Println("[info] stream created")

	go func() {
		for i := range 10 {
			s.notify(i)
			time.Sleep(200 * time.Millisecond)
		}
		s.Close()
	}()

	for v := range s.Range() {
		fmt.Println("[info] read:", v)
	}
	fmt.Println("finish")
}
