package console

import (
	"fmt"
	"os"
)

func LogError(prefix string, err error) {
	fmt.Fprintln(os.Stderr, FormatText("ERROR", Red), fmt.Sprintf("%s:", prefix), err)
}

func LogFatal(prefix string, err error) {
	LogError(prefix, err)
	os.Exit(1)
}
