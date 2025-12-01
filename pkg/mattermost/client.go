package mattermost

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	WebhookURL string
}

type Message struct {
	Text string `json:"text"`
}

func New(webhookURL string) *Client {
	return &Client{WebhookURL: webhookURL}
}

func (c *Client) Send(text string) error {
	body, _ := json.Marshal(Message{Text: text})

	req, err := http.NewRequest("POST", c.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	return nil
}
