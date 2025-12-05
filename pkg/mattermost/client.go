package mattermost

import (
	"context"
	"time"

	"innotech/pkg/http_client"
)

// Client provides methods for sending messages to Mattermost webhooks.
type Client struct {
	httpClient *http_client.Client
	webhookURL string
}

// Message represents a Mattermost webhook message.
type Message struct {
	Text string `json:"text"`
}

// New creates a new Mattermost client instance.
func New(webhookURL string) *Client {
	client := http_client.New("", 10*time.Second)
	client.SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient: client,
		webhookURL: webhookURL,
	}
}

// Send sends a text message to the Mattermost webhook.
func (c *Client) Send(text string) error {
	return c.SendWithContext(context.Background(), text)
}

// SendWithContext sends a text message to the Mattermost webhook with context.
func (c *Client) SendWithContext(ctx context.Context, text string) error {
	message := Message{Text: text}

	_, err := c.httpClient.Post(ctx, c.webhookURL, message, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendMessage sends a custom message to the Mattermost webhook.
func (c *Client) SendMessage(ctx context.Context, message Message) error {
	_, err := c.httpClient.Post(ctx, c.webhookURL, message, nil)
	if err != nil {
		return err
	}

	return nil
}
