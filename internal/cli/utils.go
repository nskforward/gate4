package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/nskforward/gate4/pkg/dir"
	"golang.org/x/term"
)

func printHeader(header string) {
	fmt.Println("=====================================================================")
	fmt.Println(" GATE4:", strings.ToUpper(header))
	fmt.Println("=====================================================================")
}

func AskPath(prompt string) string {
	return dir.Normalize(Ask(prompt))
}

func Ask(prompt string) string {
	fmt.Printf("%s: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	text = bytes.TrimSpace(text)
	return string(text)
}

func AskSecret(prompt string) string {
	fmt.Printf("%s: ", prompt)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadBytes('\r')
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	if len(text) > 0 && text[len(text)-1] == '\r' {
		text = text[:len(text)-1]
	}

	fmt.Printf("[%d bytes]\n", len(text))

	return string(text)
}
