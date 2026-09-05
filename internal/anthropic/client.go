// Package anthropic is a minimal client for the Anthropic Messages API,
// scoped to the single call this project needs: force the model to respond
// via one tool call and return that call's raw input. It has no knowledge
// of job postings, candidate profiles, or prompts — those are owned by
// internal/analyzer.
package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.anthropic.com"
	defaultModel   = "claude-sonnet-4-6"
	apiVersion     = "2023-06-01"
	maxTokens      = 2048
)

// Client talks to the Anthropic Messages API.
type Client struct {
	APIKey     string
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

// NewClient returns a Client configured with production defaults.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		BaseURL:    defaultBaseURL,
		Model:      defaultModel,
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// ToolSpec describes a single tool the model must call.
type ToolSpec struct {
	Name        string
	Description string
	InputSchema map[string]any
}

// CallTool sends system and userMessage to the Messages API, forcing the
// model to respond by calling tool, and returns that tool call's raw JSON
// input. It returns an error if the request fails, the API responds with a
// non-2xx status, or the response does not contain a call to tool.
func (c *Client) CallTool(system, userMessage string, tool ToolSpec) (json.RawMessage, error) {
	reqBody := messageRequest{
		Model:     c.Model,
		MaxTokens: maxTokens,
		System:    system,
		Messages: []message{
			{Role: "user", Content: userMessage},
		},
		Tools: []toolDef{
			{Name: tool.Name, Description: tool.Description, InputSchema: tool.InputSchema},
		},
		ToolChoice: &toolChoice{Type: "tool", Name: tool.Name},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("anthropic: encoding request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.BaseURL+"/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("anthropic: building request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("anthropic: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr apiErrorResponse
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("anthropic: API error (%d): %s", resp.StatusCode, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("anthropic: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var msgResp messageResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return nil, fmt.Errorf("anthropic: decoding response: %w", err)
	}

	for _, block := range msgResp.Content {
		if block.Type == "tool_use" && block.Name == tool.Name {
			return block.Input, nil
		}
	}

	return nil, fmt.Errorf("anthropic: response did not include a %q tool call", tool.Name)
}
