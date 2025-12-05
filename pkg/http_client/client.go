// Package http_client provides a universal HTTP client for making HTTP requests.
package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"innotech/pkg/logger"
)

// Client provides methods for making HTTP requests.
type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
}

// RequestOptions contains options for HTTP requests.
type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	Timeout time.Duration
}

// Response contains the HTTP response data.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// New creates a new HTTP client instance.
func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		headers: make(map[string]string),
	}
}

// SetHeader sets a default header for all requests.
func (c *Client) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
}

// SetHeaders sets multiple default headers for all requests.
func (c *Client) SetHeaders(headers map[string]string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	for k, v := range headers {
		c.headers[k] = v
	}
}

// Do performs an HTTP request with the given options.
func (c *Client) Do(ctx context.Context, opts RequestOptions) (*Response, error) {
	url := opts.URL
	if c.baseURL != "" && !isAbsoluteURL(opts.URL) {
		url = c.baseURL + opts.URL
	}

	var body io.Reader
	if opts.Body != nil {
		bodyBytes, err := json.Marshal(opts.Body)
		if err != nil {
			logger.Error("failed to marshal request body",
				"url", url,
				"error", err,
			)
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, url, body)
	if err != nil {
		logger.Error("failed to create HTTP request",
			"method", opts.Method,
			"url", url,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set request-specific headers
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	// Set Content-Type if body is present and not already set
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply timeout if specified
	if opts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	logger.Debug("making HTTP request",
		"method", opts.Method,
		"url", url,
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("HTTP request failed",
			"method", opts.Method,
			"url", url,
			"error", err,
		)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Warn("failed to close response body",
				"url", url,
				"error", closeErr,
			)
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read response body",
			"url", url,
			"status_code", resp.StatusCode,
			"error", err,
		)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		logger.Warn("HTTP request returned error status",
			"method", opts.Method,
			"url", url,
			"status_code", resp.StatusCode,
			"response_body", string(bodyBytes),
		)
		return &Response{
			StatusCode: resp.StatusCode,
			Body:       bodyBytes,
			Headers:    resp.Header,
		}, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	logger.Debug("HTTP request completed successfully",
		"method", opts.Method,
		"url", url,
		"status_code", resp.StatusCode,
	)

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
		Headers:    resp.Header,
	}, nil
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(ctx, RequestOptions{
		Method:  http.MethodPost,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(ctx, RequestOptions{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
	})
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(ctx, RequestOptions{
		Method:  http.MethodPut,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(ctx, RequestOptions{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
	})
}

// isAbsoluteURL checks if the URL is absolute.
func isAbsoluteURL(url string) bool {
	return len(url) > 0 && (url[0] == '/' || url[:7] == "http://" || url[:8] == "https://")
}
