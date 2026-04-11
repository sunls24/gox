package gox

import "math/rand/v2"

func RandStr(length int) string {
	if length <= 0 {
		return ""
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func PickRandom[T any](list []T) T {
	if len(list) == 0 {
		var zero T
		return zero
	}
	return list[rand.IntN(len(list))]
}
