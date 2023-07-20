package util

func Keys[K comparable, T any](m map[K]T) []K {
	if m == nil {
		return nil
	}

	keys := make([]K, len(m))
	c := 0
	for k := range m {
		keys[c] = k
		c++
	}

	return keys
}

func AnyInMap[K comparable, T any](m map[K]T, filter func(key K, val T) bool) bool {
	for key, val := range m {
		if filter(key, val) {
			return true
		}
	}

	return false
}
