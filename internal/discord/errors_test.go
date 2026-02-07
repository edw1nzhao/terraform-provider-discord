package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

// ---------- TestIsNotFound ----------

func TestIsNotFound_WithDiscordAPIError_404(t *testing.T) {
	t.Parallel()

	err := &DiscordAPIError{
		HTTPStatus: 404,
		Code:       10003,
		Message:    "Unknown Channel",
	}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true for 404 DiscordAPIError")
	}
}

func TestIsNotFound_WithDiscordAPIError_403(t *testing.T) {
	t.Parallel()

	err := &DiscordAPIError{
		HTTPStatus: 403,
		Code:       50013,
		Message:    "Missing Permissions",
	}
	if IsNotFound(err) {
		t.Error("expected IsNotFound to return false for 403 DiscordAPIError")
	}
}

func TestIsNotFound_WithOtherError(t *testing.T) {
	t.Parallel()

	err := errors.New("some random error")
	if IsNotFound(err) {
		t.Error("expected IsNotFound to return false for non-DiscordAPIError")
	}
}

func TestIsNotFound_WithNil(t *testing.T) {
	t.Parallel()

	if IsNotFound(nil) {
		t.Error("expected IsNotFound to return false for nil")
	}
}

func TestIsNotFound_WithWrappedError(t *testing.T) {
	t.Parallel()

	inner := &DiscordAPIError{
		HTTPStatus: 404,
		Code:       10003,
		Message:    "Unknown Channel",
	}
	wrapped := fmt.Errorf("wrapped: %w", inner)
	if !IsNotFound(wrapped) {
		t.Error("expected IsNotFound to return true for wrapped 404 DiscordAPIError")
	}
}

// ---------- TestDiscordAPIError_Error ----------

func TestDiscordAPIError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      DiscordAPIError
		expected string
	}{
		{
			name: "without validation errors",
			err: DiscordAPIError{
				HTTPStatus: 404,
				Code:       10003,
				Message:    "Unknown Channel",
			},
			expected: "discord API error (HTTP 404, code 10003): Unknown Channel",
		},
		{
			name: "with validation errors",
			err: DiscordAPIError{
				HTTPStatus: 400,
				Code:       50035,
				Message:    "Invalid Form Body",
				Errors:     json.RawMessage(`{"name": {"_errors": [{"code": "BASE_TYPE_REQUIRED", "message": "required"}]}}`),
			},
			expected: `discord API error (HTTP 400, code 50035): Invalid Form Body - {"name": {"_errors": [{"code": "BASE_TYPE_REQUIRED", "message": "required"}]}}`,
		},
		{
			name: "zero code",
			err: DiscordAPIError{
				HTTPStatus: 500,
				Code:       0,
				Message:    "Internal Server Error",
			},
			expected: "discord API error (HTTP 500, code 0): Internal Server Error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.err.Error()
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// ---------- TestRateLimitError_Error ----------

func TestRateLimitError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      RateLimitError
		expected string
	}{
		{
			name: "non-global rate limit",
			err: RateLimitError{
				RetryAfter: 1.5,
				Global:     false,
				Message:    "You are being rate limited.",
			},
			expected: "discord rate limited, retry after 1.50s",
		},
		{
			name: "global rate limit",
			err: RateLimitError{
				RetryAfter: 5.0,
				Global:     true,
				Message:    "You are being rate limited.",
			},
			expected: "discord rate limited (global), retry after 5.00s",
		},
		{
			name: "fractional retry after",
			err: RateLimitError{
				RetryAfter: 0.123,
				Global:     false,
				Message:    "You are being rate limited.",
			},
			expected: "discord rate limited, retry after 0.12s",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.err.Error()
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// ---------- TestDiscordAPIError_ImplementsError ----------

func TestDiscordAPIError_ImplementsError(t *testing.T) {
	t.Parallel()

	var _ error = &DiscordAPIError{}
}

// ---------- TestRateLimitError_ImplementsError ----------

func TestRateLimitError_ImplementsError(t *testing.T) {
	t.Parallel()

	var _ error = &RateLimitError{}
}
