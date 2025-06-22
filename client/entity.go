package client

type Pair[T any] struct {
	Key   string
	Value T
}

func NewPair[T any](key string, value T) Pair[T] {
	return Pair[T]{key, value}
}
