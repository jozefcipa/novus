package shared

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

func MapKeys[K string, V interface{}](obj map[K]V) []K {
	keys := make([]K, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	return keys
}

func MergeMaps[V any, M ~map[string]V](m1 M, m2 M) M {
	merged := make(M)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
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
