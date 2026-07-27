package mail

import (
	"context"
	"errors"
	stdmail "net/mail"
	"strings"
	"testing"
)

func TestNewRequiresSMTPConfig(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("New() error = nil, want configuration error")
	}
}

func TestBuildMessage(t *testing.T) {
	from := &stdmail.Address{Address: "sender@example.com"}
	to := &stdmail.Address{Address: "user@example.com"}
	message, err := buildMessage(from, to, "注册验证码", "验证码：123456")
	if err != nil {
		t.Fatalf("buildMessage() error = %v", err)
	}
	content := string(message)
	for _, want := range []string{
		"From: <sender@example.com>",
		"To: <user@example.com>",
		"Subject: =?UTF-8?q?",
		"Content-Type: text/plain; charset=UTF-8",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("message does not contain %q", want)
		}
	}
}

func TestSendCanceledContext(t *testing.T) {
	client, err := New(Config{
		Host:     "smtp.example.com",
		Port:     465,
		Username: "sender@example.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = client.Send(ctx, "user@example.com", "subject", "text")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Send() error = %v, want context.Canceled", err)
	}
}
