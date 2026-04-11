package openai

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/sunls24/gox/network/client"
	"github.com/sunls24/gox/network/header"
)

type OpenAI struct {
	baseURL string
	apiKey  string
	client  *client.Client
}

type Option func(*OpenAI)

func New(baseURL, apiKey string, opts ...Option) *OpenAI {
	oai := &OpenAI{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  strings.TrimSpace(apiKey),
		client:  client.New(),
	}

	for _, opt := range opts {
		opt(oai)
	}

	return oai
}

func WithClient(c *http.Client) Option {
	return func(oai *OpenAI) {
		oai.client = client.New(client.WithClient(c))
	}
}

func (oai *OpenAI) Responses(ctx context.Context, req Request) (io.ReadCloser, error) {
	const path = "/responses"
	return oai.postReader(ctx, path, req)
}

func (oai *OpenAI) ChatCompletions(ctx context.Context, req ChatRequest) (io.ReadCloser, error) {
	const path = "/chat/completions"
	return oai.postReader(ctx, path, req)
}

func (oai *OpenAI) postReader(ctx context.Context, path string, body any) (io.ReadCloser, error) {
	return oai.client.PostReader(
		ctx,
		oai.baseURL+path,
		body,
		header.New().ContentTypeJSON().Authorization(oai.apiKey).Get()...,
	)
}
