package nadago

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
}

type Option func(c *Client)

func NewClient(baseurl string, opts ...Option) *Client {
	client := &Client{
		apiURL:     baseurl,
		httpClient: http.DefaultClient,
	}

	for _, o := range opts {
		o(client)
	}

	return client
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithURLAuth(username, password string) Option {
	return func(c *Client) {
		// append auth parameters to URL query
		c.apiURL += fmt.Sprintf("?username=%s&password=%s", url.QueryEscape(username), url.QueryEscape(password))
	}
}
