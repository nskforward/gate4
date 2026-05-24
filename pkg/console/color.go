package console

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
)
