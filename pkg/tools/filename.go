package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func NormalizeFilename(in string) (string, error) {
	if in == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	if in[0] == '/' {
		return in, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, in), nil
}
