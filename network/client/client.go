package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sunls24/gox"
	"github.com/sunls24/gox/types"
)

var def = &http.Client{
	Timeout: time.Minute * 10,
}

func SetTimeout(timeout time.Duration) {
	def.Timeout = timeout
}

func Get(ctx context.Context, url string, header ...types.Pair[string]) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return Do(req, header...)
}

func Post(ctx context.Context, url string, body any, header ...types.Pair[string]) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, NewBody(body))
	if err != nil {
		return nil, err
	}
	return Do(req, header...)
}

func PostReader(ctx context.Context, url string, body any, header ...types.Pair[string]) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, NewBody(body))
	if err != nil {
		return nil, err
	}
	return DoReader(req, header...)
}

func Close(reader io.ReadCloser) {
	_, _ = io.Copy(io.Discard, reader)
	_ = reader.Close()
}

func Do(req *http.Request, header ...types.Pair[string]) ([]byte, error) {
	reader, err := DoReader(req, header...)
	if err != nil {
		return nil, err
	}
	defer Close(reader)
	return io.ReadAll(reader)
}

// DoReader reader need close
func DoReader(req *http.Request, header ...types.Pair[string]) (io.ReadCloser, error) {
	for _, p := range header {
		req.Header.Add(p.Key, p.Value)
	}
	resp, err := def.Do(req)
	if err != nil {
		return nil, err
	}
	if !statusOK(resp.StatusCode) {
		defer Close(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if len(body) != 0 {
			return nil, fmt.Errorf("%s: %s", resp.Status, gox.Bytes2Str(body))
		}
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func statusOK(code int) bool {
	return http.StatusOK <= code && code < http.StatusBadRequest
}

func NewBody(body any) io.Reader {
	if body == nil {
		return nil
	}
	data, _ := json.Marshal(body)
	return bytes.NewReader(data)
}
