package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sunls24/gox/types"
)

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
