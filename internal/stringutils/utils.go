package stringutils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToKebabCase(input string) string {
	return strings.ReplaceAll(
		strings.ToLower(
			strings.Join(
				strings.Fields(input),
				" ",
			),
		),
		" ",
		"-",
	)
}

func LowerFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}

	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}

	return string(lc) + s[size:]
}
