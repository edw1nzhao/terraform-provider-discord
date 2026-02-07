package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// BaseURL is the base URL for the Discord API v10.
	BaseURL = "https://discord.com/api/v10"

	// maxRetries is the maximum number of retries for failed requests.
	maxRetries = 3

	// baseBackoff is the base backoff duration for retries.
	baseBackoff = 1 * time.Second
)

// rateLimitBucket tracks rate limit state for a specific route.
type rateLimitBucket struct {
	mu        sync.Mutex
	remaining int
	resetAt   time.Time
}

// Client is a REST client for the Discord API v10.
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
	userAgent  string

	mu      sync.RWMutex
	buckets map[string]*rateLimitBucket
}

// NewClient creates a new Discord API client with the given bot token and provider version.
func NewClient(token, version string) *Client {
	ua := fmt.Sprintf("DiscordBot (terraform-provider-discord, %s)", version)
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token:     token,
		baseURL:   BaseURL,
		userAgent: ua,
		buckets:   make(map[string]*rateLimitBucket),
	}
}

// getBucket returns the rate limit bucket for a given route, creating one if it does not exist.
func (c *Client) getBucket(route string) *rateLimitBucket {
	c.mu.RLock()
	b, ok := c.buckets[route]
	c.mu.RUnlock()
	if ok {
		return b
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	// Double-check after acquiring write lock.
	if b, ok = c.buckets[route]; ok {
		return b
	}
	b = &rateLimitBucket{
		remaining: 1, // Assume we can make at least one request.
	}
	c.buckets[route] = b
	return b
}

// waitForRateLimit blocks until the rate limit for the given route bucket has reset.
func (c *Client) waitForRateLimit(ctx context.Context, bucket *rateLimitBucket) error {
	bucket.mu.Lock()
	remaining := bucket.remaining
	resetAt := bucket.resetAt
	bucket.mu.Unlock()

	if remaining <= 0 && time.Now().Before(resetAt) {
		delay := time.Until(resetAt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil
}

// updateRateLimit updates the rate limit bucket from response headers.
func (c *Client) updateRateLimit(bucket *rateLimitBucket, resp *http.Response) {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		if val, err := strconv.Atoi(remaining); err == nil {
			bucket.remaining = val
		}
	}

	if resetAfter := resp.Header.Get("X-RateLimit-Reset-After"); resetAfter != "" {
		if val, err := strconv.ParseFloat(resetAfter, 64); err == nil {
			bucket.resetAt = time.Now().Add(time.Duration(val*1000) * time.Millisecond)
		}
	}
}

// doRequest performs an HTTP request and decodes the JSON response into result.
// It handles rate limiting, retries, and error parsing.
func (c *Client) doRequest(ctx context.Context, method, route string, body interface{}, result interface{}) error {
	return c.doRequestInternal(ctx, method, route, body, result, false)
}

// doRequestNoContent performs an HTTP request that expects no response body (204).
func (c *Client) doRequestNoContent(ctx context.Context, method, route string, body interface{}) error {
	return c.doRequestInternal(ctx, method, route, body, nil, true)
}

// doRequestInternal is the core HTTP request handler with retries and rate limiting.
func (c *Client) doRequestInternal(ctx context.Context, method, route string, body interface{}, result interface{}, noContent bool) error {
	bucket := c.getBucket(route)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := baseBackoff * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Wait for any active rate limit before making the request.
		if err := c.waitForRateLimit(ctx, bucket); err != nil {
			return err
		}

		// Build request body.
		var reqBody io.Reader
		if body != nil {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewReader(jsonData)
		}

		// Build the HTTP request.
		url := c.baseURL + route
		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bot "+c.token)
		req.Header.Set("User-Agent", c.userAgent)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// Execute the request.
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}

		// Update rate limit state from response headers.
		c.updateRateLimit(bucket, resp)

		// Read response body.
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		// Handle rate limiting (429).
		if resp.StatusCode == http.StatusTooManyRequests {
			var rlErr RateLimitError
			if jsonErr := json.Unmarshal(respBody, &rlErr); jsonErr == nil {
				// Update bucket with retry-after.
				bucket.mu.Lock()
				bucket.remaining = 0
				bucket.resetAt = time.Now().Add(time.Duration(rlErr.RetryAfter*1000) * time.Millisecond)
				bucket.mu.Unlock()
			}
			lastErr = &RateLimitError{
				RetryAfter: rlErr.RetryAfter,
				Global:     rlErr.Global,
				Message:    rlErr.Message,
			}
			continue
		}

		// Handle server errors (5xx) - retry.
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("discord server error: HTTP %d", resp.StatusCode)
			continue
		}

		// Handle 404 Not Found.
		if resp.StatusCode == http.StatusNotFound {
			var apiErr DiscordAPIError
			if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil {
				apiErr.HTTPStatus = resp.StatusCode
				return &apiErr
			}
			return &DiscordAPIError{
				HTTPStatus: resp.StatusCode,
				Code:       0,
				Message:    "resource not found",
			}
		}

		// Handle other client errors (4xx).
		if resp.StatusCode >= 400 {
			var apiErr DiscordAPIError
			if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil {
				apiErr.HTTPStatus = resp.StatusCode
				return &apiErr
			}
			return &DiscordAPIError{
				HTTPStatus: resp.StatusCode,
				Code:       0,
				Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody)),
			}
		}

		// Handle successful no-content responses.
		if noContent || resp.StatusCode == http.StatusNoContent {
			return nil
		}

		// Decode the JSON response.
		if result != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}

		return nil
	}

	return fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}
