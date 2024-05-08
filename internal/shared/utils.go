package shared

import "strings"

// difference returns the elements in `a` that aren't in `b`.
// https://stackoverflow.com/a/45428032/4480179
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

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
