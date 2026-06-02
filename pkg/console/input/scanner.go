package input

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nskforward/gate4/pkg/console/output"
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

func (s *Scanner) Confirm(ctx context.Context, prompt string) bool {
	result, _ := s.ScanBool(ctx, output.FormatText(prompt, output.Yellow), nil, nil)
	return result
}

func (s *Scanner) ScanTime(ctx context.Context, prompt, layout string, defaultValue time.Time, dst *time.Time) (time.Time, error) {
	dstStr := ""
	if dst != nil {
		dstStr = (*dst).Format(layout)
	}

	defStr := ""
	if !defaultValue.IsZero() {
		defStr = defaultValue.Format(layout)
	} else {
		defStr = time.Now().Format(layout)
	}

	input, err := s.Scan(ctx, prompt, defStr, &dstStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(layout, input)
}

func (s *Scanner) ScanBool(ctx context.Context, prompt string, defaultValue, dst *bool) (bool, error) {
	prompt = fmt.Sprintf("%s (y/n)", prompt)

	dstStr := ""
	if dst != nil {
		dstStr = bool2str(*dst)
	}

	defStr := ""
	if defaultValue != nil {
		defStr = bool2str(*defaultValue)
	}

	input, err := s.Scan(ctx, prompt, defStr, &dstStr)
	if err != nil {
		return false, err
	}
	return str2bool(input), nil
}

func (s *Scanner) Scan(ctx context.Context, prompt, defaultValue string, dst *string) (string, error) {
	prompt = fmt.Sprintf("- %s:", prompt)

	if dst != nil && *dst != "" {
		fmt.Println(prompt, *dst)
		return *dst, nil
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	fmt.Print(prompt)
	fmt.Print(" ")
	buf := []byte(defaultValue)
	if len(buf) > 0 {
		fmt.Print(string(buf))
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		oneByte := make([]byte, 1)
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
				lineCh <- string(buf)
				return
			case 0x03: // Ctrl+C
				cancel() // отменяем контекст
				os.Stdout.Write([]byte("^C\r\n"))
				return // горутина завершается
			case 0x7f, 0x08: // Backspace
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
					os.Stdout.Write([]byte("\b \b"))
				}
			default:
				if c >= 0x20 {
					buf = append(buf, c)
					os.Stdout.Write([]byte{c})
				}
			}
		}
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

func (s *Scanner) ScanPassword(ctx context.Context, prompt string) (string, error) {

	prompt = fmt.Sprintf("- %s: ", prompt)

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	fmt.Print(prompt)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	passCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		var buf []byte
		oneByte := make([]byte, 1)

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
			case '\r', '\n':
				os.Stdout.Write([]byte("\r\n"))
				passCh <- string(buf)
				return
			case 0x03: // Ctrl+C
				cancel()
				os.Stdout.Write([]byte("^C\r\n"))
				return
			case 0x7f, 0x08: // Backspace
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
					updateCount()
				}
			default:
				if c >= 0x20 {
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
