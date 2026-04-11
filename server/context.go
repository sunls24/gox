package server

import (
	"context"

	"github.com/labstack/echo/v5"
)

func (srv *Server) ContextValue(key, value any) {
	srv.contextValue[key] = value
}

func (srv *Server) NewValueContext() context.Context {
	ctx := context.Background()
	for k, v := range srv.contextValue {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

func (srv *Server) contextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Request().Context()
		for k, v := range srv.contextValue {
			ctx = context.WithValue(ctx, k, v)
		}
		ctx = context.WithValue(ctx, echoContextKey{}, c)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

type echoContextKey struct {
}

func EchoContext(ctx context.Context) *echo.Context {
	return ctx.Value(echoContextKey{}).(*echo.Context)
}
