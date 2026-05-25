package console

import "strings"

func FindArg(key string, args []string) (string, bool) {
	found := false
	for _, arg := range args {
		if found {
			return arg, found
		}
		if strings.EqualFold(key, arg) {
			found = true
		}
	}
	return "", found
}
