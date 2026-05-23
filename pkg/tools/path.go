package tools

import (
	"os"
	"path/filepath"
)

func Path(in string) string {
	if in == "" {
		return in
	}
	if in[0] == '/' {
		return in
	}
	wd, err := os.Getwd()
	if err != nil {
		return in
	}
	return filepath.Join(wd, in)
}
