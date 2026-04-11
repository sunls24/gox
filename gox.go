package gox

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime/debug"
)

func MD5(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}

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

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(fmt.Sprintf("SafeGo panic: %v\n%s", err, string(debug.Stack())))
			}
		}()
		fn()
	}()
}

type Result[T any] struct {
	Value T
	Err   error
}

func Async[T any](fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)
	go func() {
		defer close(ch)
		v, err := fn()
		ch <- Result[T]{
			Value: v,
			Err:   err,
		}
	}()
	return ch
}
