package openai

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

type Request struct {
	Model        string     `json:"model"`
	Instructions string     `json:"instructions"`
	Input        []*Message `json:"input"`
	Temperature  float64    `json:"temperature"`
	Stream       bool       `json:"stream"`
	Store        bool       `json:"store"`
}
