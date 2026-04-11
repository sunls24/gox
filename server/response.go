package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v5"
)

func writeReply(c *echo.Context, r *Reply) error {
	if r == nil {
		r = OK(nil)
	}
	if r.handled {
		return nil
	}
	statusCode := statusOrOK(r.statusCode)
	if r.noBody || c.Request().Method == http.MethodHead {
		return c.NoContent(statusCode)
	}
	return c.JSON(statusCode, r.body)
}

func writeError(c *echo.Context, err error) error {
	if isCanceled(err) {
		return nil
	}
	var e *Error
	if !errors.As(err, &e) {
		return err
	}
	if e.err != nil && !isCanceled(e.err) {
		c.Logger().Error(fmt.Sprintf("%s: %v\n%s", e.message, e.err, string(debug.Stack())))
	}
	statusCode := statusOrOK(e.statusCode)
	if c.Request().Method == http.MethodHead {
		return c.NoContent(statusCode)
	}
	return c.JSON(statusCode, e.Envelope())
}

func statusOrOK(statusCode int) int {
	if statusCode == 0 {
		return http.StatusOK
	}
	return statusCode
}

func isCanceled(err error) bool {
	return errors.Is(err, context.Canceled)
}

func httpErrorHandler(errorHandler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		if isCanceled(err) {
			return
		}
		if code := echo.StatusCode(err); code > 0 && code < http.StatusInternalServerError {
			errorHandler(c, err)
			return
		}
		c.Logger().Error(fmt.Sprintf("%v\n%s", err, string(debug.Stack())))
		errorHandler(c, echo.NewHTTPError(http.StatusInternalServerError, "未知异常"))
	}
}
