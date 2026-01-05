package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Empty struct{}

const unknownError = "未知异常"

// Wrap T must struct
func Wrap[T any, R any](fn func(ctx context.Context, req T) (R, error)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusOK, ErrMsg("参数解析异常"))
		}
		resp, err := fn(c.Request().Context(), req)
		if err != nil {
			r := new(Response)
			if errors.As(err, &r) {
				if r.Err != nil {
					c.Logger().Error(r.Message, slog.Any("err", r.Err))
				}
				return c.JSON(http.StatusOK, r)
			}
			return err
		}
		return c.JSON(http.StatusOK, data(resp))
	}
}

func Start(address string, init func(*echo.Echo)) error {
	e := echo.New()
	errorHandler := echo.DefaultHTTPErrorHandler(false)
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		c.Logger().Error(err.Error())
		errorHandler(c, echo.NewHTTPError(http.StatusInternalServerError, unknownError))
	}
	e.Use(contextMiddleware)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll: true,
	}))
	init(e)
	return e.Start(address)
}
