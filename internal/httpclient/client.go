package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RequestResult represents the result of a single HTTP request
type RequestResult struct {
	ResponseStatus  int
	ResponseTimeMs  int64
	ResponseHeaders map[string]string
	Success         bool
	Error           string
}

// Client wraps HTTP client with timeout control
type Client struct {
	httpClient *http.Client
	timeout    time.Duration
}

// New creates a new HTTP client with specified timeout in milliseconds
func New(timeoutMs int) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		},
		timeout: time.Duration(timeoutMs) * time.Millisecond,
	}
}

// DoRequest performs an HTTP request and returns the result
func (c *Client) DoRequest(method, url string, headers map[string]string, body string) *RequestResult {
	start := time.Now()

	// Build request
	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return &RequestResult{
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Success:        false,
			Error:          fmt.Sprintf("failed to create request: %v", err),
		}
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Execute request with context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &RequestResult{
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Success:        false,
			Error:          fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, _ := io.ReadAll(resp.Body)

	// Collect response headers
	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	// Determine success (2xx status codes)
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	result := &RequestResult{
		ResponseStatus:  resp.StatusCode,
		ResponseTimeMs:  time.Since(start).Milliseconds(),
		ResponseHeaders: respHeaders,
		Success:         success,
	}

	if !success {
		result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return result
}
