package di

import (
	"log/slog"
	"os"

	"github.com/nskforward/gate4/pkg/console"
)

func initLogger() *slog.Logger {
	return console.NewLogger(os.Stdout, slog.LevelDebug)
}
