package maputils

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
