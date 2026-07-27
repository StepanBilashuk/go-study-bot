// Package claude is a tiny HTTP client for the Anthropic Messages API plus the
// strict parsers for the JSON shapes the bot asks Claude to return.
//
// It deliberately uses net/http rather than the official SDK: the project's
// hard constraint is minimal dependencies (spec §1), and this is "an HTTP
// client for the Anthropic API" with no transitive dependencies.
package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	messagesURL      = "https://api.anthropic.com/v1/messages"
	anthropicVersion = "2023-06-01"
	maxTokens        = 8192 // generous: adaptive thinking + a small JSON reply
)

// Client talks to the Anthropic Messages API.
type Client struct {
	apiKey string
	model  string
	http   *http.Client
}

// New returns a client for the given API key and model id.
func New(apiKey, model string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{Timeout: 120 * time.Second},
	}
}

type apiRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	System    string       `json:"system,omitempty"`
	Messages  []apiMessage `json:"messages"`
}

type apiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type apiResponse struct {
	Content    []apiBlock `json:"content"`
	StopReason string     `json:"stop_reason"`
	Error      *apiError  `json:"error"`
}

type apiBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type apiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Complete sends a single-turn request (system prompt + one user message) and
// returns the concatenated text of the reply. Thinking blocks, if any, are
// ignored — only text blocks are joined.
func (c *Client) Complete(ctx context.Context, system, user string) (string, error) {
	body, err := json.Marshal(apiRequest{
		Model:     c.model,
		MaxTokens: maxTokens,
		System:    system,
		Messages:  []apiMessage{{Role: "user", Content: user}},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, messagesURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("call anthropic: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var parsed apiResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decode response (status %d): %w", resp.StatusCode, err)
	}
	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil {
			return "", fmt.Errorf("anthropic %s: %s", parsed.Error.Type, parsed.Error.Message)
		}
		return "", fmt.Errorf("anthropic http %d: %s", resp.StatusCode, string(raw))
	}
	if parsed.StopReason == "refusal" {
		return "", fmt.Errorf("anthropic declined the request")
	}

	var sb strings.Builder
	for _, b := range parsed.Content {
		if b.Type == "text" {
			sb.WriteString(b.Text)
		}
	}
	text := strings.TrimSpace(sb.String())
	if text == "" {
		return "", fmt.Errorf("empty response from anthropic (stop_reason %q)", parsed.StopReason)
	}
	return text, nil
}

// CompleteJSON runs Complete, extracts the JSON payload, and hands it to parse.
// parse must both unmarshal and validate. On a parse failure it retries once
// with a clarifying instruction before returning the error (spec §14 Step 3).
func (c *Client) CompleteJSON(ctx context.Context, system, user string, parse func([]byte) error) error {
	raw, err := c.Complete(ctx, system, user)
	if err != nil {
		return err
	}
	if err := parse([]byte(ExtractJSON(raw))); err == nil {
		return nil
	}

	retry := user + "\n\nYour previous reply was not valid JSON matching the schema. " +
		"Reply with ONLY the JSON value — no prose, no explanation, no markdown code fences."
	raw, err = c.Complete(ctx, system, retry)
	if err != nil {
		return err
	}
	if err := parse([]byte(ExtractJSON(raw))); err != nil {
		return fmt.Errorf("claude did not return valid JSON after one retry: %w", err)
	}
	return nil
}

// ExtractJSON strips a surrounding ```json ... ``` fence (or a leading/trailing
// stray fence) and returns the inner text, so a fenced reply still parses.
func ExtractJSON(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimPrefix(s, "json")
	s = strings.TrimPrefix(s, "JSON")
	if i := strings.LastIndex(s, "```"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}
