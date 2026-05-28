package console

import "fmt"

type Color struct {
	R int
	G int
	B int
}

var (
	Red     = Color{255, 0, 0}
	Green   = Color{0, 255, 0}
	Blue    = Color{0, 0, 255}
	Yellow  = Color{255, 255, 0}
	Cyan    = Color{0, 255, 255}
	Magenta = Color{255, 0, 255}
	White   = Color{255, 255, 255}
	Black   = Color{0, 0, 0}
	Gray100 = Color{100, 100, 100}
	Gray200 = Color{200, 200, 200}
)

func (color Color) String() string {
	return fmt.Sprintf("38;5;%d", rgbTo256(color.R, color.G, color.B))
}

func rgbTo256(r, g, b int) int {
	r6 := r * 5 / 255
	g6 := g * 5 / 255
	b6 := b * 5 / 255
	return 16 + 36*r6 + 6*g6 + b6
}
