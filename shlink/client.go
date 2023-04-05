package shlink

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc

	apiKey string
	url    *url.URL

	http http.Client
}

func NewClient(ctx context.Context, apiKey string, apiUrl string) (*Client, error) {
	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	client := &Client{apiKey: apiKey, url: parsedUrl, ctx: ctx}
	if ctx == nil {
		client.ctx = context.Background()
	}
	client.ctx, client.cancel = context.WithCancel(client.ctx)
	return client, nil
}

func (c *Client) doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(c.ctx, method, c.url.JoinPath(url).String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", c.apiKey)
	return c.http.Do(req)
}
