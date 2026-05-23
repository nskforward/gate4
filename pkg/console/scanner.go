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
// ScanPassword читает пароль без эхо-вывода, отображая счётчик символов.
// prompt — строка, которая будет выведена перед началом ввода (например, "- secret: ").
func (s *Scanner) ScanPassword(ctx context.Context, prompt string) (string, error) {

	prompt = fmt.Sprintf("- %s: ", prompt)

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	// Выводим приглашение один раз перед входом в raw‑режим.
	fmt.Print(prompt)

	passCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		var buf []byte
		oneByte := make([]byte, 1)

		// Обновляет строку: возвращает каретку, печатает prompt + счётчик, стирает хвост.
		updateCount := func() {
			fmt.Fprintf(os.Stdout, "\r%s[%d chars]\033[K", prompt, len(buf))
		}

		for {
			_, err := os.Stdin.Read(oneByte)
			if err != nil {
				errCh <- err
				return
			}
			c := oneByte[0]

			switch c {
			case '\r', '\n': // Enter
				os.Stdout.Write([]byte("\r\n"))
				passCh <- string(buf)
				return
			case 0x7f, 0x08: // Backspace
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
					updateCount()
				}
			default:
				if c >= 0x20 { // печатные символы
					buf = append(buf, c)
					updateCount()
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case pass := <-passCh:
		return pass, nil
	}
}
