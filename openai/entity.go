package openai

import "encoding/json"

type Role string

const (
	RUser      Role = "user"
	RSystem    Role = "system"
	RAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func NewMessage(role Role, content string) *Message {
	return &Message{Role: role, Content: content}
}

func UserMessage(content string) *Message {
	return NewMessage(RUser, content)
}

func SystemMessage(content string) *Message {
	return NewMessage(RSystem, content)
}

func AssistantMessage(content string) *Message {
	return NewMessage(RAssistant, content)
}

type Reasoning struct {
	Effort string `json:"effort,omitempty"` // minimal/low/medium/high/max; none 表示关闭思考模式
}

type Request struct {
	Model        string     `json:"model"`
	Instructions string     `json:"instructions,omitempty"`
	Input        Input      `json:"input"`
	Reasoning    *Reasoning `json:"reasoning,omitempty"`
	Temperature  *float64   `json:"temperature,omitempty"`
	Stream       *bool      `json:"stream,omitempty"`
	Store        *bool      `json:"store,omitempty"`
}

type ChatRequest struct {
	Model       string     `json:"model"`
	Messages    []*Message `json:"messages"`
	Temperature *float64   `json:"temperature,omitempty"`
	Stream      *bool      `json:"stream,omitempty"`
}

type Input struct {
	text     *string
	messages []*Message
}

func StringInput(text string) Input {
	return Input{text: &text}
}

func MessagesInput(messages ...*Message) Input {
	return Input{messages: append([]*Message{}, messages...)}
}

func (i Input) MarshalJSON() ([]byte, error) {
	if i.text != nil {
		return json.Marshal(*i.text)
	}
	if i.messages != nil {
		return json.Marshal(i.messages)
	}
	return []byte("null"), nil
}
