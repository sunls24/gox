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
)

var def = &http.Client{
	Timeout: time.Minute * 10,
}

func SetTimeout(timeout time.Duration) {
	def.Timeout = timeout
}

func Get(ctx context.Context, url string, header ...Pair[string]) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if len(header) > 0 {
		for _, p := range header {
			req.Header.Add(p.Key, p.Value)
		}
	}
	return Do(req)
}

func Post(ctx context.Context, url string, body any, header ...Pair[string]) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, NewBody(body))
	if err != nil {
		return nil, err
	}
	if len(header) > 0 {
		for _, p := range header {
			req.Header.Add(p.Key, p.Value)
		}
	}
	return Do(req)
}

func Do(req *http.Request) ([]byte, error) {
	resp, err := def.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if !statusOK(resp.StatusCode) {
		if len(body) != 0 {
			return nil, errors.New(fmt.Sprintf("%s: %s", resp.Status, gox.Bytes2Str(body)))
		}
		return nil, errors.New(resp.Status)
	}
	return body, nil
}

func DoStream(req *http.Request) (io.ReadCloser, error) {
	resp, err := def.Do(req)
	if err != nil {
		return nil, err
	}
	if !statusOK(resp.StatusCode) {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
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
