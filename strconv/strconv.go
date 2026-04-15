package strconv

import "strings"

func EscapeString(s string) string {
	var bob strings.Builder
	bob.Grow(len(s))
	for _, b := range []byte(s) {
		if b == '\\' || b == '"' {
			bob.WriteByte('\\')
		}
		bob.WriteByte(b)
	}
	return bob.String()
}
