package server

import (
	"context"

	"github.com/labstack/echo/v5"
)

type Empty struct{}

// Wrap T must struct
func Wrap[T any, R any](fn func(ctx context.Context, req T) (R, error)) echo.HandlerFunc {
	return wrap(true, func(ctx context.Context, req T) (*Reply, error) {
		resp, err := fn(ctx, req)
		if err != nil {
			return nil, err
		}
		return OK(resp), nil
	})
}

func WrapEmpty(fn func(ctx context.Context) error) echo.HandlerFunc {
	return wrap(false, func(ctx context.Context, req Empty) (*Reply, error) {
		if err := fn(ctx); err != nil {
			return nil, err
		}
		return OK(nil), nil
	})
}

func WrapReq[T any](fn func(ctx context.Context, req T) error) echo.HandlerFunc {
	return wrap(true, func(ctx context.Context, req T) (*Reply, error) {
		if err := fn(ctx, req); err != nil {
			return nil, err
		}
		return OK(nil), nil
	})
}

func WrapResp[R any](fn func(ctx context.Context) (R, error)) echo.HandlerFunc {
	return wrap(false, func(ctx context.Context, req Empty) (*Reply, error) {
		resp, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		return OK(resp), nil
	})
}

func WrapReply[T any](fn func(ctx context.Context, req T) (*Reply, error)) echo.HandlerFunc {
	return wrap(true, fn)
}

func WrapReplyResp(fn func(ctx context.Context) (*Reply, error)) echo.HandlerFunc {
	return wrap(false, func(ctx context.Context, req Empty) (*Reply, error) {
		return fn(ctx)
	})
}

func wrap[T any](bind bool, fn func(ctx context.Context, req T) (*Reply, error)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req T
		if bind {
			if err := c.Bind(&req); err != nil {
				return writeError(c, BadParam())
			}
		}

		reply, err := fn(c.Request().Context(), req)
		if err != nil {
			return writeError(c, err)
		}
		return writeReply(c, reply)
	}
}
