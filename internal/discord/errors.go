package discord

import (
	"encoding/json"
	"errors"
	"fmt"
)

// DiscordAPIError represents an error response from the Discord API.
type DiscordAPIError struct {
	HTTPStatus int             `json:"-"`
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Errors     json.RawMessage `json:"errors,omitempty"`
}

// Error implements the error interface.
func (e *DiscordAPIError) Error() string {
	if e.Errors != nil {
		return fmt.Sprintf("discord API error (HTTP %d, code %d): %s - %s", e.HTTPStatus, e.Code, e.Message, string(e.Errors))
	}
	return fmt.Sprintf("discord API error (HTTP %d, code %d): %s", e.HTTPStatus, e.Code, e.Message)
}

// RateLimitError represents a 429 rate limit response from the Discord API.
type RateLimitError struct {
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
	Message    string  `json:"message"`
}

// Error implements the error interface.
func (e *RateLimitError) Error() string {
	if e.Global {
		return fmt.Sprintf("discord rate limited (global), retry after %.2fs", e.RetryAfter)
	}
	return fmt.Sprintf("discord rate limited, retry after %.2fs", e.RetryAfter)
}

// IsNotFound returns true if the error is a DiscordAPIError with HTTP status 404.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *DiscordAPIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatus == 404
	}

	return false
}
