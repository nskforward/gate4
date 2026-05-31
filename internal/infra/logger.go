package infra

import (
	"log/slog"
	"os"

	"github.com/nskforward/gate4/pkg/console"
)

func InitLogger() {
	slog.SetDefault(console.NewLogger(os.Stdout, slog.LevelDebug))
}
