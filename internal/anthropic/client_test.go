package anthropic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &Client{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		Model:      defaultModel,
		HTTPClient: server.Client(),
	}
}

func testTool() ToolSpec {
	return ToolSpec{
		Name:        "submit_fit_analysis",
		Description: "Submit the fit analysis",
		InputSchema: map[string]any{"type": "object"},
	}
}

func TestCallTool_Success(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"content": [
				{"type": "tool_use", "name": "submit_fit_analysis", "input": {"match_score": 72}}
			]
		}`))
	})

	input, err := client.CallTool("system prompt", "job text", testTool())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed struct {
		MatchScore int `json:"match_score"`
	}
	if err := json.Unmarshal(input, &parsed); err != nil {
		t.Fatalf("failed to unmarshal tool input: %v", err)
	}
	if parsed.MatchScore != 72 {
		t.Errorf("expected match_score 72, got %d", parsed.MatchScore)
	}
}

func TestCallTool_HTTPError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": {"type": "api_error", "message": "server exploded"}}`))
	})

	_, err := client.CallTool("system", "job text", testTool())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "server exploded") {
		t.Errorf("expected error to mention API message, got: %v", err)
	}
}

func TestCallTool_MalformedResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	})

	_, err := client.CallTool("system", "job text", testTool())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCallTool_MissingToolCall(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"content": [{"type": "text", "text": "I refuse to use tools"}]}`))
	})

	_, err := client.CallTool("system", "job text", testTool())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "did not include") {
		t.Errorf("expected error about missing tool call, got: %v", err)
	}
}
