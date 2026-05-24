package console

import (
	"fmt"
	"strings"
)

func FormatText(text string, color Color, styles ...Style) string {
	colorCode := fmt.Sprintf("38;5;%d", rgbTo256(color.R, color.G, color.B))
	codes := make([]string, 0, len(styles)+1)
	for _, s := range styles {
		codes = append(codes, fmt.Sprint(int(s)))
	}
	codes = append(codes, colorCode)
	prefix := "\033[" + strings.Join(codes, ";") + "m"
	suffix := "\033[0m"
	return prefix + text + suffix
}

func rgbTo256(r, g, b int) int {
	r6 := r * 5 / 255
	g6 := g * 5 / 255
	b6 := b * 5 / 255
	return 16 + 36*r6 + 6*g6 + b6
}
