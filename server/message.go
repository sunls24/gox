package server

import (
	"errors"
	"fmt"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`

	err        error
	statusCode int
	skipW      bool
}

func SkipW() *Response {
	return &Response{
		skipW: true,
	}
}

func (r *Response) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

func ErrMsg(msg string) *Response {
	return &Response{
		Code:       -1,
		Message:    msg,
		statusCode: http.StatusOK,
	}
}

func ErrMsgf(format string, a ...any) *Response {
	return ErrMsg(fmt.Sprintf(format, a...))
}

func (r *Response) WithErr(err error) *Response {
	if err == nil {
		return r
	}
	//goland:noinspection GoTypeAssertionOnErrors
	if child, ok := err.(*Response); ok {
		r.Message = fmt.Sprintf("%s: %s", r.Message, child.Message)
		r.err = errors.Join(r.err, child.err)
	} else {
		r.err = err
	}
	return r
}

func (r *Response) WithStatusCode(statusCode int) *Response {
	r.statusCode = statusCode
	return r
}

func Data(data any) *Response {
	return &Response{
		Message:    "ok",
		Data:       data,
		statusCode: http.StatusOK,
	}
}

func BadParam() *Response {
	return ErrMsg("请求参数不符合要求")
}
