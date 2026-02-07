package discord

import (
	"encoding/json"
	"testing"
)

// ---------- TestSnowflake_String ----------

func TestSnowflake_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		s        Snowflake
		expected string
	}{
		{name: "normal ID", s: Snowflake("123456789012345678"), expected: "123456789012345678"},
		{name: "empty", s: Snowflake(""), expected: ""},
		{name: "zero", s: Snowflake("0"), expected: "0"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.s.String()
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// ---------- TestSnowflake_IsEmpty ----------

func TestSnowflake_IsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		s        Snowflake
		expected bool
	}{
		{name: "empty string", s: Snowflake(""), expected: true},
		{name: "zero string", s: Snowflake("0"), expected: true},
		{name: "valid ID", s: Snowflake("123456789012345678"), expected: false},
		{name: "non-zero numeric", s: Snowflake("1"), expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.s.IsEmpty()
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// ---------- TestSnowflake_MarshalJSON ----------

func TestSnowflake_MarshalJSON(t *testing.T) {
	t.Parallel()

	// Snowflake is a string alias, so JSON marshaling should produce a quoted string.
	type wrapper struct {
		ID Snowflake `json:"id"`
	}

	tests := []struct {
		name     string
		input    wrapper
		expected string
	}{
		{
			name:     "normal snowflake",
			input:    wrapper{ID: Snowflake("123456789012345678")},
			expected: `{"id":"123456789012345678"}`,
		},
		{
			name:     "empty snowflake",
			input:    wrapper{ID: Snowflake("")},
			expected: `{"id":""}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := string(data)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// ---------- TestSnowflake_UnmarshalJSON ----------

func TestSnowflake_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type wrapper struct {
		ID Snowflake `json:"id"`
	}

	tests := []struct {
		name     string
		input    string
		expected Snowflake
	}{
		{
			name:     "string value",
			input:    `{"id":"123456789012345678"}`,
			expected: Snowflake("123456789012345678"),
		},
		{
			name:     "empty string",
			input:    `{"id":""}`,
			expected: Snowflake(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var w wrapper
			err := json.Unmarshal([]byte(tc.input), &w)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if w.ID != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, w.ID)
			}
		})
	}
}

// ---------- TestSnowflake_JSONRoundTrip ----------

func TestSnowflake_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	type wrapper struct {
		ID Snowflake `json:"id"`
	}

	original := wrapper{ID: Snowflake("987654321098765432")}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded wrapper
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("round-trip failed: expected %q, got %q", original.ID, decoded.ID)
	}
}

// ---------- TestSnowflake_InStruct ----------

func TestSnowflake_InStruct(t *testing.T) {
	t.Parallel()

	// Test that a Discord-like JSON response unmarshals correctly.
	jsonData := `{"id":"111222333444555666","name":"general","guild_id":"999888777666555444"}`

	var ch Channel
	err := json.Unmarshal([]byte(jsonData), &ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch.ID != Snowflake("111222333444555666") {
		t.Errorf("expected ID %q, got %q", "111222333444555666", ch.ID)
	}
	if ch.GuildID == nil || *ch.GuildID != Snowflake("999888777666555444") {
		t.Error("expected GuildID to be 999888777666555444")
	}
}
