package dir

import (
	"os"
	"path/filepath"
)

func Normalize(path string) string {
	if path != "" && path[0] == '/' {
		return path
	}
	wd, err := os.Getwd()
	if err != nil {
		return path
	}
	return filepath.Join(wd, path)
}
