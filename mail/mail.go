package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	stdmail "net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Timeout  time.Duration
}

type Client struct {
	config Config
	from   *stdmail.Address
}

func New(config Config) (*Client, error) {
	config.Host = strings.TrimSpace(config.Host)
	config.Username = strings.TrimSpace(config.Username)
	config.From = strings.TrimSpace(config.From)
	if config.From == "" {
		config.From = config.Username
	}
	if config.Timeout <= 0 {
		config.Timeout = defaultTimeout
	}
	if config.Host == "" {
		return nil, errors.New("mail host is required")
	}
	if config.Port < 1 || config.Port > 65535 {
		return nil, errors.New("mail port is invalid")
	}
	if config.Username == "" || config.Password == "" {
		return nil, errors.New("mail username and password are required")
	}
	from, err := stdmail.ParseAddress(config.From)
	if err != nil {
		return nil, fmt.Errorf("parse mail from address: %w", err)
	}
	return &Client{config: config, from: from}, nil
}

func (c *Client) Send(ctx context.Context, to, subject, text string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.ContainsAny(subject, "\r\n") {
		return errors.New("mail subject contains a newline")
	}
	recipient, err := stdmail.ParseAddress(strings.TrimSpace(to))
	if err != nil {
		return fmt.Errorf("parse mail recipient: %w", err)
	}
	message, err := buildMessage(c.from, recipient, subject, text)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()
	address := net.JoinHostPort(c.config.Host, strconv.Itoa(c.config.Port))
	connection, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return smtpError(ctx, "dial SMTP server", err)
	}
	defer connection.Close()
	stopClose := context.AfterFunc(ctx, func() {
		_ = connection.Close()
	})
	defer stopClose()
	if deadline, ok := ctx.Deadline(); ok {
		if err := connection.SetDeadline(deadline); err != nil {
			return smtpError(ctx, "set SMTP deadline", err)
		}
	}

	secureConnection := tls.Client(connection, &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: c.config.Host,
	})
	if err := secureConnection.HandshakeContext(ctx); err != nil {
		return smtpError(ctx, "handshake SMTP TLS", err)
	}
	client, err := smtp.NewClient(secureConnection, c.config.Host)
	if err != nil {
		return smtpError(ctx, "create SMTP client", err)
	}
	defer client.Close()

	if err := client.Auth(smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)); err != nil {
		return smtpError(ctx, "authenticate SMTP", err)
	}
	if err := client.Mail(c.from.Address); err != nil {
		return smtpError(ctx, "set SMTP sender", err)
	}
	if err := client.Rcpt(recipient.Address); err != nil {
		return smtpError(ctx, "set SMTP recipient", err)
	}
	writer, err := client.Data()
	if err != nil {
		return smtpError(ctx, "start SMTP message", err)
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return smtpError(ctx, "write SMTP message", err)
	}
	if err := writer.Close(); err != nil {
		return smtpError(ctx, "close SMTP message", err)
	}
	_ = client.Quit()
	return nil
}

func smtpError(ctx context.Context, operation string, err error) error {
	if contextErr := ctx.Err(); contextErr != nil {
		return contextErr
	}
	return fmt.Errorf("%s: %w", operation, err)
}

func buildMessage(from, to *stdmail.Address, subject, text string) ([]byte, error) {
	var message bytes.Buffer
	fmt.Fprintf(&message, "From: %s\r\n", from.String())
	fmt.Fprintf(&message, "To: %s\r\n", to.String())
	fmt.Fprintf(&message, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject))
	fmt.Fprintf(&message, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	writer := quotedprintable.NewWriter(&message)
	if _, err := writer.Write([]byte(text)); err != nil {
		return nil, fmt.Errorf("encode mail text: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close mail text encoder: %w", err)
	}
	return message.Bytes(), nil
}
