package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/sunls24/gox/types"
)

type Client struct {
	client *http.Client
}

func New() *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Minute * 10,
		},
	}
}

func NewWithJar() *Client {
	c := New()
	c.client.Jar, _ = cookiejar.New(nil)
	return c
}

var def = New()

func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

func SetTimeout(timeout time.Duration) {
	def.SetTimeout(timeout)
}

func (c *Client) Get(ctx context.Context, url string, header ...types.Pair[string]) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, header...)
}

func Get(ctx context.Context, url string, header ...types.Pair[string]) ([]byte, error) {
	return def.Get(ctx, url, header...)
}

func (c *Client) Post(ctx context.Context, url string, body any, header ...types.Pair[string]) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, NewBody(body))
	if err != nil {
		return nil, err
	}
	return c.Do(req, header...)
}

func Post(ctx context.Context, url string, body any, header ...types.Pair[string]) ([]byte, error) {
	return def.Post(ctx, url, body, header...)
}

func (c *Client) PostReader(ctx context.Context, url string, body any, header ...types.Pair[string]) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, NewBody(body))
	if err != nil {
		return nil, err
	}
	return c.DoReader(req, header...)
}

func PostReader(ctx context.Context, url string, body any, header ...types.Pair[string]) (io.ReadCloser, error) {
	return def.PostReader(ctx, url, body, header...)
}

func Close(reader io.ReadCloser) {
	_, _ = io.Copy(io.Discard, reader)
	_ = reader.Close()
}

func (c *Client) Do(req *http.Request, header ...types.Pair[string]) ([]byte, error) {
	reader, err := c.DoReader(req, header...)
	if err != nil {
		return nil, err
	}
	defer Close(reader)
	return io.ReadAll(reader)
}

func Do(req *http.Request, header ...types.Pair[string]) ([]byte, error) {
	return def.Do(req, header...)
}

// DoReader reader need close
func (c *Client) DoReader(req *http.Request, header ...types.Pair[string]) (io.ReadCloser, error) {
	for _, p := range header {
		req.Header.Set(p.Key, p.Value)
	}
	resp, err := c.client.Do(req)
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
			return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
		}
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func DoReader(req *http.Request, header ...types.Pair[string]) (io.ReadCloser, error) {
	return def.DoReader(req, header...)
}

func statusOK(code int) bool {
	return http.StatusOK <= code && code < http.StatusMultipleChoices
}

func NewBody(body any) io.Reader {
	if body == nil {
		return nil
	}
	data, _ := json.Marshal(body)
	return bytes.NewReader(data)
}
