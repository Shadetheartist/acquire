package util

func RemoveAt[T comparable](slice []T, idx int) []T {
	if slice == nil {
		return slice
	}

	if len(slice) == 0 {
		return slice
	}

	slice = append(slice[:idx], slice[idx+1:]...)

	return slice
}

func IndexOf[T comparable](slice []T, val T) (int, bool) {
	if slice == nil {
		return 0, false
	}

	if len(slice) == 0 {
		return 0, false
	}

	for idx, v := range slice {
		if v == val {
			return idx, true
		}
	}

	return 0, false
}

func UniqueElements[T comparable](slice []T) []T {
	encountered := make(map[T]bool)
	result := make([]T, 0)

	for _, item := range slice {
		if !encountered[item] {
			encountered[item] = true
			result = append(result, item)
		}
	}

	return result
}

func Map[T comparable, T2 any](slice []T, mapper func(val T) T2) []T2 {
	result := make([]T2, len(slice))

	for i, item := range slice {
		result[i] = mapper(item)
	}

	return result
}

func Filter[T comparable](slice []T, filter func(val T) bool) []T {
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if filter(item) {
			result = append(result, item)
		}
	}

	return result
}

func Find[T comparable](slice []T, filter func(val T) bool) (int, T) {
	var z T

	for idx, item := range slice {
		if filter(item) {
			return idx, item
		}
	}

	return -1, z
}

func Any[T comparable](slice []T, filter func(val T) bool) bool {
	for _, item := range slice {
		if filter(item) {
			return true
		}
	}

	return false
}

func Clone[T comparable](slice []T) []T {
	clone := make([]T, len(slice))
	copy(clone, slice)
	return clone
}

func Pop[T comparable](slice []T) (T, []T) {
	var val T

	if slice == nil {
		return val, nil
	}

	if len(slice) == 0 {
		return val, slice
	}

	val, slice = slice[0], slice[1:]
	return val, slice
}
