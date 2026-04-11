package server

import (
	"errors"
	"fmt"
	"net/http"
)

type Envelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	code    int
	message string

	err        error
	statusCode int
}

func ErrMsg(msg string) *Error {
	return &Error{
		code:       -1,
		message:    msg,
		statusCode: http.StatusOK,
	}
}

func ErrMsgf(format string, a ...any) *Error {
	return ErrMsg(fmt.Sprintf(format, a...))
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) WithErr(err error) *Error {
	if err == nil {
		return e
	}
	if child, ok := err.(*Error); ok {
		e.message = fmt.Sprintf("%s: %s", e.message, child.message)
		e.err = joinErr(e.err, child.err)
	} else {
		e.err = joinErr(e.err, err)
	}
	return e
}

func joinErr(base, err error) error {
	if err == nil {
		return base
	}
	if base == nil {
		return err
	}
	return errors.Join(base, err)
}

func (e *Error) WithStatusCode(statusCode int) *Error {
	e.statusCode = statusCode
	return e
}

func (e *Error) Envelope() Envelope {
	return Envelope{
		Code:    e.code,
		Message: e.message,
	}
}

func BadParam() *Error {
	return ErrMsg("请求参数不符合要求")
}

type Reply struct {
	statusCode int
	body       any
	noBody     bool
	handled    bool
}

func OK(data any) *Reply {
	return &Reply{
		statusCode: http.StatusOK,
		body: Envelope{
			Message: "ok",
			Data:    data,
		},
	}
}

func StatusCode(statusCode int) *Reply {
	return &Reply{
		statusCode: statusCode,
		noBody:     true,
	}
}

func Handled() *Reply {
	return &Reply{
		handled: true,
	}
}

func (r *Reply) WithStatusCode(statusCode int) *Reply {
	r.statusCode = statusCode
	return r
}
