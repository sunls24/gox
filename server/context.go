package server

import (
	"context"

	"github.com/labstack/echo/v5"
)

func (s *Server) ContextValue(key, value any) {
	s.contextValues[key] = value
}

func (s *Server) NewValueContext() context.Context {
	return s.withContextValues(context.Background())
}

func (s *Server) contextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := s.withContextValues(c.Request().Context())
		ctx = context.WithValue(ctx, echoContextKey{}, c)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func (s *Server) withContextValues(ctx context.Context) context.Context {
	for k, v := range s.contextValues {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

type echoContextKey struct {
}

func EchoContext(ctx context.Context) *echo.Context {
	return ctx.Value(echoContextKey{}).(*echo.Context)
}
