package console

import "fmt"

type textColor string

const (
	reset  textColor = "\033[0m"
	red    textColor = "\033[31m"
	green  textColor = "\033[32m"
	normal textColor = "\033[39m"
)

func RedText(s string) string {
	return sprintColor(red, s)
}

func GreenText(s string) string {
	return sprintColor(green, s)
}

func sprintColor(code textColor, s string) string {
	return fmt.Sprint(code, s, reset)
}
