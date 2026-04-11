package client

import (
	"context"
	"io"
	"net/http"

	"github.com/sunls24/gox/types"
)

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
