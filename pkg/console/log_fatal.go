package console

import (
	"log/slog"
	"os"
)

func LogFatal(prefix string, err error) {
	slog.Error(prefix, "error", err.Error())
	os.Exit(1)
}
