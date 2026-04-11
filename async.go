package gox

import (
	"fmt"
	"log/slog"
	"runtime/debug"
)

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
