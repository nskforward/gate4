package console

import (
	"fmt"
	"os"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func LogDebug(msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format(timeFormat), "DEBUG", msg)
}

func LogWarn(msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format(timeFormat), FormatText("WARN ", Yellow, Bold), FormatText(msg, White, Bold))
}

func LogInfo(msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format(timeFormat), FormatText("INFO ", Green, Bold), FormatText(msg, White, Bold))
}

func LogError(prefix string, err error) {
	fmt.Fprintln(os.Stderr, time.Now().Format(timeFormat), FormatText("ERROR", Red, Bold), fmt.Sprintf("%s:", FormatText(prefix, White, Bold)), err)
}

func LogFatal(prefix string, err error) {
	LogError(prefix, err)
	os.Exit(1)
}
