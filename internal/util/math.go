package util

func Max[T numeric](a T, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T numeric](a T, b T) T {
	if a < b {
		return a
	}

	return b
}
