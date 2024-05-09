package shared

import "strings"

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
