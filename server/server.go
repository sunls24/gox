package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Empty struct{}

// Wrap T must struct
func Wrap[T any, R any](fn func(ctx context.Context, req T) (R, error)) echo.HandlerFunc {
	return wrap(true, fn)
}

func WrapEmpty(fn func(ctx context.Context) error) echo.HandlerFunc {
	return wrap(false, func(ctx context.Context, req Empty) (*Empty, error) {
		return nil, fn(ctx)
	})
}

func WrapReq[T any](fn func(ctx context.Context, req T) error) echo.HandlerFunc {
	return wrap(true, func(ctx context.Context, req T) (*Empty, error) {
		return nil, fn(ctx, req)
	})
}

func WrapResp[R any](fn func(ctx context.Context) (R, error)) echo.HandlerFunc {
	return wrap(false, func(ctx context.Context, req Empty) (R, error) {
		return fn(ctx)
	})
}

func wrap[T any, R any](bind bool, fn func(ctx context.Context, req T) (R, error)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req T
		if bind {
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusOK, BadParam())
			}
		}

		resp, err := fn(c.Request().Context(), req)
		if err == nil {
			return c.JSON(http.StatusOK, Data(resp))
		}
		//goland:noinspection GoTypeAssertionOnErrors
		if r, ok := err.(*Response); ok {
			if r.err != nil {
				c.Logger().Error(fmt.Sprintf("%s: %v\n%s", r.Message, r.err, string(debug.Stack())))
			}
			if r.skipW {
				return nil
			}
			return c.JSON(r.statusCode, r)
		}
		return err
	}
}

type Server struct {
	Echo         *echo.Echo
	contextValue map[any]any
}

func Start(address string, init func(*Server)) error {
	e := echo.New()
	errorHandler := echo.DefaultHTTPErrorHandler(false)
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		c.Logger().Error(fmt.Sprintf("%v\n%s", err, string(debug.Stack())))
		errorHandler(c, echo.NewHTTPError(http.StatusInternalServerError, "未知异常"))
	}
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll: true,
	}))
	srv := &Server{
		contextValue: make(map[any]any),
		Echo:         e,
	}
	init(srv)
	e.Use(srv.contextMiddleware)
	return e.Start(address)
}
