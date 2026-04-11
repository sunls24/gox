package gox

func Ptr[T any](value T) *T {
	return &value
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
