package client

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	client *http.Client
}

type Option func(*Client)

func WithClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.client = client
		}
	}
}

func WithJar() Option {
	return func(c *Client) {
		c.client.Jar, _ = cookiejar.New(nil)
	}
}

func New(opts ...Option) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: time.Minute * 10,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

var def = New()

func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

func SetTimeout(timeout time.Duration) {
	def.SetTimeout(timeout)
}
