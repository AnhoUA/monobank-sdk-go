package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BaseURL is the default base URL for the Monobank API.
const BaseURL = "https://api.monobank.ua"

// Client is a base HTTP client for Monobank API with retry logic.
type Client struct {
	BaseURL       string
	HTTPClient    *http.Client
	RetryInterval time.Duration
	MaxRetries    int
	Token         string
}

// Option allows configuring the client.
type Option func(*Client)

// WithHTTPClient returns an Option that sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithRetryInterval returns an Option that sets the interval between retry attempts.
func WithRetryInterval(interval time.Duration) Option {
	return func(c *Client) {
		c.RetryInterval = interval
	}
}

// NewClient creates a new base client with default settings.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		Token:         token,
		BaseURL:       BaseURL,
		RetryInterval: 5 * time.Second,
		MaxRetries:    3,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ErrorResponse represents an error returned by the Monobank API.
type ErrorResponse struct {
	StatusCode       int    `json:"-"`
	ErrorDescription string `json:"errorDescription"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.ErrorDescription)
}

// IsTooManyRequests checks if the error is a result of exceeding request limits.
func IsTooManyRequests(err error) bool {
	if err == nil {
		return false
	}
	var errResp *ErrorResponse
	if errors.As(err, &errResp) {
		return errResp.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// Do performs an HTTP request with retry logic for 429 errors.
func (c *Client) Do(ctx context.Context, method, path string, body any, target any, auth bool) error {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	var lastErr error
	for attempt := range c.MaxRetries + 1 {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.RetryInterval):
			}
		}

		err := c.execute(ctx, method, path, bodyBytes, target, auth)
		if err == nil {
			return nil
		}

		lastErr = err
		if !IsTooManyRequests(err) {
			return err
		}
	}

	return lastErr
}

func (c *Client) execute(ctx context.Context, method, path string, body []byte, target any, auth bool) error {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	if auth && c.Token != "" {
		req.Header.Set("X-Token", c.Token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		errResp.StatusCode = resp.StatusCode
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("http error: %d", resp.StatusCode)
		}
		return &errResp
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}
