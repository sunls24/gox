package gox

import (
	"math/rand/v2"
)

func In[T comparable](target T, list ...T) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func PickRandom[T any](list []T) T {
	if len(list) == 0 {
		var zero T
		return zero
	}
	return list[rand.IntN(len(list))]
}

func Map[T any, R any](input []T, fn func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](input []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(input))
	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
