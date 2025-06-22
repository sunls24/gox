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

func SimplePrompt(user string, system string) []Message {
	var list []Message
	if system != "" {
		list = append(list, SystemMessage(system))
	}
	list = append(list, UserMessage(user))
	return list
}

type StreamFunc func(data []byte) error

type ReqChat struct {
	Stream      bool      `json:"stream"`
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`

	OnStream StreamFunc `json:"-"`
	OnStart  func()     `json:"-"`
}

func (rc *ReqChat) check() error {
	if rc.Stream && rc.OnStream == nil {
		return errors.New("OnStream is required when stream is true")
	}
	if len(rc.Messages) == 0 {
		return errors.New("request messages is empty")
	}
	return nil
}
