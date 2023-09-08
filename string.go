package ignition

import "strings"

func CapitalizeString(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if len(s) == 0 {
		return ""
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}
