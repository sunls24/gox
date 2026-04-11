package header

import (
	"fmt"
	"slices"

	"github.com/sunls24/gox/types"
)

var contentTypeJSON = types.NewPair("Content-Type", "application/json")

var chromeHeaders = []types.Pair[string]{
	types.NewPair("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8"),
	types.NewPair("Priority", "u=0, i"),
	types.NewPair("Sec-CH-UA", "\"Google Chrome\";v=\"147\", \"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"147\""),
	types.NewPair("Sec-CH-UA-Mobile", "?0"),
	types.NewPair("Sec-CH-UA-Platform", "\"macOS\""),
	types.NewPair("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"),
}

var chromeDocumentHeaders = append(
	slices.Clone(chromeHeaders),
	types.NewPair("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"),
	types.NewPair("Sec-Fetch-Dest", "document"),
	types.NewPair("Sec-Fetch-Mode", "navigate"),
	types.NewPair("Sec-Fetch-Site", "none"),
	types.NewPair("Sec-Fetch-User", "?1"),
	types.NewPair("Upgrade-Insecure-Requests", "1"),
)

type Builder struct {
	headers []types.Pair[string]
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Add(headers ...types.Pair[string]) *Builder {
	b.headers = append(b.headers, headers...)
	return b
}

func (b *Builder) ChromeHeaders() *Builder {
	return b.Add(slices.Clone(chromeHeaders)...)
}

func (b *Builder) ChromeDocumentHeaders() *Builder {
	return b.Add(slices.Clone(chromeDocumentHeaders)...)
}

func (b *Builder) ContentTypeJSON() *Builder {
	return b.Add(contentTypeJSON)
}

func (b *Builder) Authorization(token string) *Builder {
	return b.Add(types.NewPair("Authorization", fmt.Sprintf("Bearer %s", token)))
}

func (b *Builder) Referer(referer string) *Builder {
	return b.Add(types.NewPair("Referer", referer))
}

func (b *Builder) Get() []types.Pair[string] {
	return b.headers
}
