package console

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/term"
)

type Scanner struct {
	closeOnce sync.Once
	done      chan struct{}
	scanErr   atomic.Pointer[error]
}

func NewScanner() *Scanner {
	return &Scanner{
		done: make(chan struct{}),
	}
}

func (s *Scanner) Close() {
	s.closeOnce.Do(func() {
		close(s.done)
	})
}

func (s *Scanner) ScanTime(ctx context.Context, layout string) (time.Time, error) {
	input, err := s.Scan(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(layout, input)
}

func (s *Scanner) ScanBool(ctx context.Context) (bool, error) {
	input, err := s.Scan(ctx)
	if err != nil {
		return false, err
	}

	input = strings.ToLower(input)

	if slices.Contains([]string{"y", "yes", "ok", "1", "true"}, input) {
		return true, nil
	}
	if slices.Contains([]string{"n", "no", "false", "0"}, input) {
		return false, nil
	}

	return false, fmt.Errorf("unrecognized bool input: %s", input)
}

// Scan читает одну строку с отображением вводимых символов.
// Не использует буферизацию, совместим с последующим вызовом ScanPassword.
func (s *Scanner) Scan(ctx context.Context) (string, error) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		var buf []byte
		b := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(b)
			if err != nil {
				errCh <- err
				return
			}
			if n == 0 {
				continue
			}
			if b[0] == '\n' || b[0] == '\r' {
				break
			}
			buf = append(buf, b[0])
		}
		lineCh <- string(buf)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case line := <-lineCh:
		return line, nil
	}
}

// ScanPassword читает строку без отображения (как пароль).
func (s *Scanner) ScanPassword(ctx context.Context) (string, error) {
	passCh := make(chan []byte, 1)
	errCh := make(chan error, 1)

	go func() {
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			errCh <- err
			return
		}
		passCh <- password
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case pass := <-passCh:
		// ReadPassword не выводит символ новой строки после ввода,
		// поэтому добавляем его сами.
		fmt.Println()
		return string(pass), nil
	}
}
