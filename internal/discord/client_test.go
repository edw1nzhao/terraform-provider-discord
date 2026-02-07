package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient creates a Client pointing at a httptest.Server for testing.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := NewClient("test-token", "test")
	client.baseURL = server.URL
	return client, server
}

// ---------- TestNewClient ----------

func TestNewClient(t *testing.T) {
	t.Parallel()

	c := NewClient("my-bot-token", "1.2.3")

	if c.baseURL != BaseURL {
		t.Errorf("expected baseURL %q, got %q", BaseURL, c.baseURL)
	}
	if c.token != "my-bot-token" {
		t.Errorf("expected token %q, got %q", "my-bot-token", c.token)
	}
	expectedUA := "DiscordBot (terraform-provider-discord, 1.2.3)"
	if c.userAgent != expectedUA {
		t.Errorf("expected userAgent %q, got %q", expectedUA, c.userAgent)
	}
	if c.buckets == nil {
		t.Error("expected buckets map to be non-nil")
	}
	if c.httpClient == nil {
		t.Error("expected httpClient to be non-nil")
	}
}

// ---------- TestDoRequest_Success ----------

func TestDoRequest_Success(t *testing.T) {
	t.Parallel()

	type payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers.
		if got := r.Header.Get("Authorization"); got != "Bot test-token" {
			t.Errorf("expected Authorization header %q, got %q", "Bot test-token", got)
		}
		if got := r.Header.Get("User-Agent"); !strings.Contains(got, "terraform-provider-discord") {
			t.Errorf("expected User-Agent to contain terraform-provider-discord, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(payload{ID: "123", Name: "test-channel"})
	})
	defer server.Close()

	var result payload
	err := client.doRequest(context.Background(), http.MethodGet, "/channels/123", nil, &result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "123" {
		t.Errorf("expected ID %q, got %q", "123", result.ID)
	}
	if result.Name != "test-channel" {
		t.Errorf("expected Name %q, got %q", "test-channel", result.Name)
	}
}

// ---------- TestDoRequest_404_ReturnsDiscordAPIError ----------

func TestDoRequest_404_ReturnsDiscordAPIError(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code": 10003, "message": "Unknown Channel"}`))
	})
	defer server.Close()

	var result map[string]interface{}
	err := client.doRequest(context.Background(), http.MethodGet, "/channels/999", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*DiscordAPIError)
	if !ok {
		t.Fatalf("expected *DiscordAPIError, got %T: %v", err, err)
	}
	if apiErr.HTTPStatus != 404 {
		t.Errorf("expected HTTPStatus 404, got %d", apiErr.HTTPStatus)
	}
	if apiErr.Code != 10003 {
		t.Errorf("expected Code 10003, got %d", apiErr.Code)
	}
	if apiErr.Message != "Unknown Channel" {
		t.Errorf("expected Message %q, got %q", "Unknown Channel", apiErr.Message)
	}
}

// ---------- TestDoRequest_RateLimit_RetriesAfterDelay ----------

func TestDoRequest_RateLimit_RetriesAfterDelay(t *testing.T) {
	t.Parallel()

	var attempt atomic.Int32

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := attempt.Add(1)
		if n == 1 {
			w.Header().Set("Retry-After", "0.01")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message": "You are being rate limited.", "retry_after": 0.01, "global": false}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "123"}`))
	})
	defer server.Close()

	// Use a very short backoff so the test doesn't take long. The retry logic
	// uses exponential backoff *in addition* to the rate-limit wait, but the
	// base backoff is 1s by default. We cannot change the constant, so this
	// test will take ~1s for the backoff sleep on retry attempt 1.
	// That is acceptable for correctness.

	type result struct {
		ID string `json:"id"`
	}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/rate-limit", nil, &res)
	if err != nil {
		t.Fatalf("expected no error after retry, got %v", err)
	}
	if res.ID != "123" {
		t.Errorf("expected ID %q, got %q", "123", res.ID)
	}
	if got := attempt.Load(); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

// ---------- TestDoRequest_ServerError_Retries ----------

func TestDoRequest_ServerError_Retries(t *testing.T) {
	t.Parallel()

	var attempt atomic.Int32

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := attempt.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`Internal Server Error`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	})
	defer server.Close()

	type result struct {
		OK bool `json:"ok"`
	}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/server-error", nil, &res)
	if err != nil {
		t.Fatalf("expected no error after retry, got %v", err)
	}
	if !res.OK {
		t.Error("expected OK to be true")
	}
	if got := attempt.Load(); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

// ---------- TestDoRequest_MaxRetries_ReturnsError ----------

func TestDoRequest_MaxRetries_ReturnsError(t *testing.T) {
	t.Parallel()

	var attempt atomic.Int32

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempt.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`Internal Server Error`))
	})
	defer server.Close()

	type result struct {
		OK bool `json:"ok"`
	}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/max-retries", nil, &res)
	if err == nil {
		t.Fatal("expected error after max retries, got nil")
	}
	if !strings.Contains(err.Error(), "request failed after") {
		t.Errorf("expected error message to contain 'request failed after', got %q", err.Error())
	}
	// maxRetries is 3, so total attempts = 1 initial + 3 retries = 4.
	if got := attempt.Load(); got != 4 {
		t.Errorf("expected 4 attempts (1 + %d retries), got %d", maxRetries, got)
	}
}

// ---------- TestDoRequest_ContextCancellation ----------

func TestDoRequest_ContextCancellation(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Should never be reached because context is already cancelled.
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	type result struct{}
	var res result
	err := client.doRequest(ctx, http.MethodGet, "/test/cancelled", nil, &res)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	if ctx.Err() == nil {
		t.Fatal("expected context error to be non-nil")
	}
}

// ---------- TestDoRequest_InvalidJSON ----------

func TestDoRequest_InvalidJSON(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json!!!`))
	})
	defer server.Close()

	type result struct {
		ID string `json:"id"`
	}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/bad-json", nil, &res)
	if err == nil {
		t.Fatal("expected error from invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode response") {
		t.Errorf("expected error message to contain 'failed to decode response', got %q", err.Error())
	}
}

// ---------- TestDoRequest_NoBody ----------

func TestDoRequest_NoBody(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := client.doRequestNoContent(context.Background(), http.MethodDelete, "/test/no-content", nil)
	if err != nil {
		t.Fatalf("expected no error for 204 No Content, got %v", err)
	}
}

// ---------- TestDoRequest_NoBody_StatusOK_NoContent ----------

func TestDoRequest_NoBody_ViaStatusCode(t *testing.T) {
	t.Parallel()

	// Test that doRequest with a 204 status code returns nil even when result is provided.
	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	type result struct {
		ID string `json:"id"`
	}
	var res result
	err := client.doRequest(context.Background(), http.MethodDelete, "/test/no-content-2", nil, &res)
	if err != nil {
		t.Fatalf("expected no error for 204 No Content, got %v", err)
	}
}

// ---------- TestDoRequest_SendsRequestBody ----------

func TestDoRequest_SendsRequestBody(t *testing.T) {
	t.Parallel()

	type reqPayload struct {
		Name string `json:"name"`
	}
	type respPayload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
		}
		var body reqPayload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(respPayload{ID: "456", Name: body.Name})
	})
	defer server.Close()

	var res respPayload
	err := client.doRequest(context.Background(), http.MethodPost, "/guilds", reqPayload{Name: "my-guild"}, &res)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ID != "456" {
		t.Errorf("expected ID %q, got %q", "456", res.ID)
	}
	if res.Name != "my-guild" {
		t.Errorf("expected Name %q, got %q", "my-guild", res.Name)
	}
}

// ---------- TestDoRequest_4xxClientError ----------

func TestDoRequest_4xxClientError(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"code": 50013, "message": "Missing Permissions"}`))
	})
	defer server.Close()

	type result struct{}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/forbidden", nil, &res)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*DiscordAPIError)
	if !ok {
		t.Fatalf("expected *DiscordAPIError, got %T", err)
	}
	if apiErr.HTTPStatus != 403 {
		t.Errorf("expected HTTPStatus 403, got %d", apiErr.HTTPStatus)
	}
	if apiErr.Code != 50013 {
		t.Errorf("expected Code 50013, got %d", apiErr.Code)
	}
	if apiErr.Message != "Missing Permissions" {
		t.Errorf("expected Message %q, got %q", "Missing Permissions", apiErr.Message)
	}
}

// ---------- TestDoRequest_UpdatesRateLimitHeaders ----------

func TestDoRequest_UpdatesRateLimitHeaders(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "4")
		w.Header().Set("X-RateLimit-Reset-After", "1.5")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer server.Close()

	type result struct {
		OK bool `json:"ok"`
	}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/rate-headers", nil, &res)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	route := "/test/rate-headers"
	bucket := client.getBucket(route)
	bucket.mu.Lock()
	remaining := bucket.remaining
	resetAt := bucket.resetAt
	bucket.mu.Unlock()

	if remaining != 4 {
		t.Errorf("expected bucket remaining 4, got %d", remaining)
	}
	// resetAt should be in the near future (within 2 seconds from now).
	if time.Until(resetAt) < 0 || time.Until(resetAt) > 3*time.Second {
		t.Errorf("expected resetAt to be ~1.5s in the future, got %v from now", time.Until(resetAt))
	}
}

// ---------- TestRateLimitBucket_Concurrent ----------

func TestRateLimitBucket_Concurrent(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.Header().Set("X-RateLimit-Remaining", "5")
		w.Header().Set("X-RateLimit-Reset-After", "0.001")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"1"}`))
	})
	defer server.Close()

	const goroutines = 10
	var wg sync.WaitGroup
	errs := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			type result struct {
				ID string `json:"id"`
			}
			var res result
			// All goroutines use the same route to share a bucket.
			errs[idx] = client.doRequest(context.Background(), http.MethodGet, "/test/concurrent", nil, &res)
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d returned error: %v", i, err)
		}
	}

	// All requests should have completed.
	if got := requestCount.Load(); got != int32(goroutines) {
		t.Errorf("expected %d requests, got %d", goroutines, got)
	}
}

// ---------- TestGetBucket_CreatesAndReuses ----------

func TestGetBucket_CreatesAndReuses(t *testing.T) {
	t.Parallel()

	client := NewClient("tok", "v")

	b1 := client.getBucket("/channels/1")
	b2 := client.getBucket("/channels/1")
	b3 := client.getBucket("/channels/2")

	if b1 != b2 {
		t.Error("expected same bucket for same route")
	}
	if b1 == b3 {
		t.Error("expected different bucket for different route")
	}
	if b1.remaining != 1 {
		t.Errorf("expected initial remaining to be 1, got %d", b1.remaining)
	}
}

// ---------- TestDoRequest_404_NonJSONBody ----------

func TestDoRequest_404_NonJSONBody(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`not found`))
	})
	defer server.Close()

	type result struct{}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/404-text", nil, &res)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*DiscordAPIError)
	if !ok {
		t.Fatalf("expected *DiscordAPIError, got %T", err)
	}
	if apiErr.HTTPStatus != 404 {
		t.Errorf("expected HTTPStatus 404, got %d", apiErr.HTTPStatus)
	}
	if apiErr.Message != "resource not found" {
		t.Errorf("expected Message %q, got %q", "resource not found", apiErr.Message)
	}
}

// ---------- TestDoRequest_4xxClientError_NonJSONBody ----------

func TestDoRequest_4xxClientError_NonJSONBody(t *testing.T) {
	t.Parallel()

	client, server := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`bad request body`))
	})
	defer server.Close()

	type result struct{}
	var res result
	err := client.doRequest(context.Background(), http.MethodGet, "/test/400-text", nil, &res)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*DiscordAPIError)
	if !ok {
		t.Fatalf("expected *DiscordAPIError, got %T", err)
	}
	if apiErr.HTTPStatus != 400 {
		t.Errorf("expected HTTPStatus 400, got %d", apiErr.HTTPStatus)
	}
	expected := fmt.Sprintf("HTTP %d: %s", 400, "bad request body")
	if apiErr.Message != expected {
		t.Errorf("expected Message %q, got %q", expected, apiErr.Message)
	}
}

// ---------- TestWaitForRateLimit_ContextCancelled ----------

func TestWaitForRateLimit_ContextCancelled(t *testing.T) {
	t.Parallel()

	client := NewClient("tok", "v")
	bucket := client.getBucket("/test/wait-cancel")

	// Set the bucket to rate-limited state with a long reset time.
	bucket.mu.Lock()
	bucket.remaining = 0
	bucket.resetAt = time.Now().Add(10 * time.Second)
	bucket.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.waitForRateLimit(ctx, bucket)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// ---------- TestWaitForRateLimit_NotLimited ----------

func TestWaitForRateLimit_NotLimited(t *testing.T) {
	t.Parallel()

	client := NewClient("tok", "v")
	bucket := client.getBucket("/test/wait-ok")

	// Bucket has remaining > 0, so should not wait.
	bucket.mu.Lock()
	bucket.remaining = 5
	bucket.mu.Unlock()

	start := time.Now()
	err := client.waitForRateLimit(context.Background(), bucket)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elapsed > 50*time.Millisecond {
		t.Errorf("expected no delay, took %v", elapsed)
	}
}
