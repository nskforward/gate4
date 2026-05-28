package console

import (
	"log/slog"
	"os"
	"testing"
)

func TestSLog(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("message 1", "a1", 1) // {"level":"INFO","msg":"message 1","a1":1}

	logger2 := logger.With("b1", 1)
	logger2.Info("message 2", "a2", 2) // {"level":"INFO","msg":"message 2","b1":1,"a2":2}

	logger3 := logger2.WithGroup("c1")
	logger3.Info("message 3", "a3", 3) // {"level":"INFO","msg":"message 3","b1":1,"c1":{"a3":3}}

	logger4 := logger3.With("b2", 2)
	logger4.Info("message 4", "a4", 4) // {"level":"INFO","msg":"message 4","b1":1,"c1":{"b2":2,"a4":4}}

	logger5 := logger4.WithGroup("c2")
	logger5.Info("message 5", "d4", 3) // {"level":"INFO","msg":"message 5","b1":1,"c1":{"b2":2,"c2":{"d4":3}}}
}
