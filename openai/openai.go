package openai

import (
	"context"
	"io"
	"strings"

	"github.com/sunls24/gox/network/client"
	"github.com/sunls24/gox/network/header"
)

type OpenAI struct {
	baseURL string
	apiKey  string
}

func New(baseURL, apiKey string) *OpenAI {
	return &OpenAI{strings.TrimRight(baseURL, "/"), strings.TrimSpace(apiKey)}
}

func (oai *OpenAI) Responses(ctx context.Context, req Request) (io.ReadCloser, error) {
	const PATH = "/responses"
	return client.PostReader(ctx, oai.baseURL+PATH, req, header.New().ContentTypeJSON().Authorization(oai.apiKey).Get()...)
}
