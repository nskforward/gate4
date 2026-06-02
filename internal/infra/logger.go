package infra

import (
	"log/slog"
	"os"

	"github.com/nskforward/gate4/pkg/console/output"
)

func InitLogger() {
	slog.SetDefault(output.NewLogger(os.Stdout, slog.LevelDebug))
}
