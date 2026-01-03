package openai

import "errors"

type Role string

const (
	RUser      Role = "user"
	RSystem    Role = "system"
	RAssistant Role = "assistant"
	RTool      Role = "tool"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func NewMessage(role Role, content string) Message {
	return Message{Role: role, Content: content}
}

func UserMessage(content string) Message {
	return NewMessage(RUser, content)
}

func SystemMessage(content string) Message {
	return NewMessage(RSystem, content)
}

func AssistantMessage(content string) Message {
	return NewMessage(RAssistant, content)
}

func StartPrompt(user string, system string) []Message {
	var list []Message
	if system != "" {
		list = append(list, SystemMessage(system))
	}
	return append(list, UserMessage(user))
}

type StreamFunc func(data []byte) error

type ReqChat struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`

	OnStream StreamFunc `json:"-"`
	OnStart  func()     `json:"-"`
}

func (rc *ReqChat) check(stream bool) error {
	if rc.Model == "" {
		return errors.New("request model is required")
	}
	if len(rc.Messages) == 0 {
		return errors.New("request messages is required")
	}
	if rc.Temperature < 0 {
		rc.Temperature = 0
	}
	rc.Stream = stream
	if stream && rc.OnStream == nil {
		return errors.New("OnStream is required when stream is true")
	}
	return nil
}
