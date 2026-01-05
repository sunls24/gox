package server

import (
	"context"

	"github.com/labstack/echo/v5"
)

type contextKey struct {
}

func contextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, contextKey{}, c)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func Context(ctx context.Context) *echo.Echo {
	return ctx.Value(contextKey{}).(*echo.Echo)
}
