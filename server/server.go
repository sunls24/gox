package server

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Server struct {
	Echo          *echo.Echo
	contextValues map[any]any
}

func New(init func(*Server)) *Server {
	e := echo.New()
	errorHandler := echo.DefaultHTTPErrorHandler(false)
	e.HTTPErrorHandler = httpErrorHandler(errorHandler)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll: true,
	}))
	srv := &Server{
		contextValues: make(map[any]any),
		Echo:          e,
	}
	e.Use(srv.contextMiddleware)
	if init != nil {
		init(srv)
	}
	return srv
}

func (s *Server) Start(address string) error {
	return s.Echo.Start(address)
}

func (s *Server) StartWithConfig(ctx context.Context, sc echo.StartConfig) error {
	return sc.Start(ctx, s.Echo)
}

func Start(address string, init func(*Server)) error {
	return New(init).Start(address)
}

func StartWithConfig(ctx context.Context, sc echo.StartConfig, init func(*Server)) error {
	return New(init).StartWithConfig(ctx, sc)
}
