package openai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/sunls24/gox"
	"github.com/sunls24/gox/client"
	"github.com/tidwall/gjson"
)

type OpenAI struct {
	baseURL string
	apiKey  string
}

func New(baseURL, apiKey string) *OpenAI {
	return &OpenAI{baseURL, apiKey}
}

func (oai *OpenAI) Chat(ctx context.Context, rc ReqChat) (string, error) {
	const path = "/v1/chat/completions"

	if err := rc.check(); err != nil {
		return "", err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, oai.baseURL+path, client.NewBody(rc))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oai.apiKey))

	if rc.Stream {
		stream, err := client.DoStream(req)
		if err != nil {
			return "", err
		}

		if rc.OnStart != nil {
			rc.OnStart()
		}
		return "", streamLoop(ctx, stream, rc.OnStream)
	}

	body, err := client.Do(req)
	if err != nil {
		return "", err
	}

	return gjson.GetBytes(body, "choices.0.message.content").String(), nil
}

var (
	ssePrefix = []byte("data: ")
	sseDone   = []byte("[DONE]")
)

func streamLoop(ctx context.Context, stream io.ReadCloser, onStream StreamFunc) error {
	const (
		thinkStart = "<think>"
		thinkEnd   = "</think>"
	)

	done := make(chan error, 1)
	data := make(chan []byte, 1)
	go fixedWrite(ctx, data, done, onStream)

	think := 0
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !bytes.HasPrefix(line, ssePrefix) {
			continue
		}
		line = line[len(ssePrefix):]
		if bytes.Equal(line, sseDone) {
			break
		}

		var content string
		if think >= 0 {
			content = gjson.GetBytes(line, "choices.0.delta.reasoning_content").String()
			if content != "" && think == 0 {
				content = thinkStart + content
				think = 1
			}
		}

		if content == "" {
			content = gjson.GetBytes(line, "choices.0.delta.content").String()
			if content == "" {
				continue
			}
			if think == 1 {
				content = thinkEnd + content
			}
			think = -1
		}

		select {
		case data <- gox.Str2Bytes(content):
		case err := <-done:
			_ = stream.Close()
			return err
		case <-ctx.Done():
			_ = stream.Close()
			return ctx.Err()
		}
	}
	_ = stream.Close()
	close(data)
	if err := <-done; err != nil {
		return err
	}
	return scanner.Err()
}

func fixedWrite(ctx context.Context, data <-chan []byte, done chan<- error, onStream StreamFunc) {
	const maxLen = 64
	const interval = time.Millisecond * 100

	defer func() {
		if err := recover(); err != nil {
			done <- fmt.Errorf("%v", err)
		}
		close(done)
	}()

	buffer := make([]byte, 0, maxLen)
	end := false

	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			done <- ctx.Err()
			return
		case b, ok := <-data:
			if !ok {
				data = nil
				end = true
				break
			}
			buffer = append(buffer, b...)
		case <-tick.C:
			if len(buffer) == 0 {
				if end {
					return
				}
				continue
			}

			cut := gox.If(len(buffer) > maxLen, maxLen, len(buffer))
			for cut > 0 && !utf8.Valid(buffer[:cut]) {
				cut--
			}
			if cut == 0 {
				continue
			}

			b := buffer[:cut]
			buffer = buffer[len(b):]

			if err := onStream(b); err != nil {
				done <- err
				return
			}
		}
	}
}
