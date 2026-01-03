package openai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
	"unicode/utf8"

	"github.com/sunls24/gox"
	"github.com/sunls24/gox/network/client"
	"github.com/sunls24/gox/network/header"
	"github.com/tidwall/gjson"
)

type OpenAI struct {
	baseURL string
	apiKey  string
}

func New(baseURL, apiKey string) *OpenAI {
	return &OpenAI{baseURL, apiKey}
}

const chatPath = "/chat/completions"

func (oai *OpenAI) Chat(ctx context.Context, rc ReqChat) (string, error) {
	if err := rc.check(false); err != nil {
		return "", err
	}
	body, err := client.Post(ctx, oai.baseURL+chatPath, rc, header.ContentTypeJson, header.Authorization(oai.apiKey))
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(body, "choices.0.message.content").String(), nil
}

func (oai *OpenAI) ChatStream(ctx context.Context, rc ReqChat) error {
	if err := rc.check(true); err != nil {
		return err
	}

	reader, err := client.PostReader(ctx, oai.baseURL+chatPath, rc, header.ContentTypeJson, header.Authorization(oai.apiKey))
	if err != nil {
		return err
	}
	if rc.OnStart != nil {
		rc.OnStart()
	}
	return streamLoop(reader, rc.OnStream)
}

var (
	ssePrefix = []byte("data: ")
	sseDone   = []byte("[DONE]")
)

func streamLoop(stream io.ReadCloser, onStream StreamFunc) error {
	const (
		thinkStart = "<think>"
		thinkEnd   = "</think>"
	)
	done := make(chan error, 1)
	data := make(chan []byte, 1)
	go fixedWrite(data, done, onStream)

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
		}
	}
	close(data)
	_ = stream.Close()
	if err := <-done; err != nil {
		return err
	}
	return scanner.Err()
}

func fixedWrite(data <-chan []byte, done chan<- error, onStream StreamFunc) {
	const maxLen = 64
	const interval = time.Millisecond * 100
	defer func() {
		if err := recover(); err != nil {
			done <- fmt.Errorf("%v", err)
		}
		close(done)
	}()

	end := false
	buffer := make([]byte, 0, maxLen)
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case b, ok := <-data:
			if !ok {
				data = nil
				end = true
				break
			}
			buffer = append(buffer, b...)
		case <-tick.C:
			cut := len(buffer)
			if cut == 0 {
				if end {
					return
				}
				continue
			}

			if cut > maxLen {
				cut = maxLen + (cut-maxLen)/4
			}
			for cut > 0 && !utf8.Valid(buffer[:cut]) {
				cut--
			}
			if cut == 0 {
				continue
			}

			b := buffer[:cut]
			buffer = buffer[cut:]
			if err := onStream(b); err != nil {
				done <- err
				return
			}
		}
	}
}
