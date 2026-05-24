package console

import "fmt"

type textColor string

const (
	reset  textColor = "\033[0m"
	red    textColor = "\033[31m"
	green  textColor = "\033[32m"
	yellow textColor = "\033[93m"
	normal textColor = "\033[39m"
)

func RedText(s string) string {
	return sprintColor(red, s)
}

func GreenText(s string) string {
	return sprintColor(green, s)
}

func YellowText(s string) string {
	return sprintColor(yellow, s)
}

func sprintColor(code textColor, s string) string {
	return fmt.Sprint(code, s, reset)
}
