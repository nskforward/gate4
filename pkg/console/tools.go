package console

import (
	"slices"
	"strings"
)

func bool2str(in bool) string {
	if in {
		return "y"
	}
	return "n"
}

func str2bool(s string) bool {
	return slices.Contains([]string{"y", "k", "ye", "ya", "yes", "ok", "okay", "1", "true", "sure", "yeah", "yep", "yup", "yah", "ja", "well", "approve", "confirm"}, strings.Trim(strings.ToLower(s), "!"))
}
