package server

import (
	"fmt"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`

	Err error `json:"-"`
}

func (r *Response) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

func ErrMsg(msg string) *Response {
	return &Response{
		Code:    -1,
		Message: msg,
	}
}

func ErrMsgf(format string, a ...any) *Response {
	return ErrMsg(fmt.Sprintf(format, a...))
}

func (r *Response) WithErr(err error) *Response {
	r.Err = err
	return r
}

func data[T any](data T) *Response {
	return &Response{
		Message: "ok",
		Data:    data,
	}
}
