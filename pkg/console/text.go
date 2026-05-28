package console

import (
	"fmt"
	"strings"
)

const resetCode = "\033[0m"

func FormatText(text string, color Color, styles ...Style) string {
	return buildPrefix(color, styles...) + text + resetCode
}

func buildPrefix(color Color, styles ...Style) string {
	codes := make([]string, 0, len(styles)+1)
	for _, s := range styles {
		codes = append(codes, fmt.Sprint(int(s)))
	}
	codes = append(codes, color.String())
	prefix := "\033[" + strings.Join(codes, ";") + "m"
	return prefix
}
