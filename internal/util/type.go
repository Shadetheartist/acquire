package util

func AsType[T any](action any) T {
	_action, ok := action.(T)
	if !ok {
		panic("action not of correct type")
	}
	return _action
}
