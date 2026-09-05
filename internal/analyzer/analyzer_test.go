package analyzer

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/jdgomez/job-assistant/internal/anthropic"
	"github.com/jdgomez/job-assistant/internal/profile"
)

type fakeCaller struct {
	response       json.RawMessage
	err            error
	capturedSystem string
	capturedUser   string
}

func (f *fakeCaller) CallTool(system, userMessage string, tool anthropic.ToolSpec) (json.RawMessage, error) {
	f.capturedSystem = system
	f.capturedUser = userMessage
	if f.err != nil {
		return nil, f.err
	}
	return f.response, nil
}

func testProfile() profile.Profile {
	return profile.Profile{
		Name:         "Jane Doe",
		Skills:       []string{"Go", "PostgreSQL"},
		Achievements: []string{"Cut CI pipeline time from 4h to 10min"},
	}
}

func TestAnalyze_Success(t *testing.T) {
	caller := &fakeCaller{
		response: json.RawMessage(`{
			"match_score": 72,
			"strengths": ["13+ years backend experience"],
			"gaps": ["No commercial Go experience"],
			"cover_letter_focus": ["Lead with CI/CD transformation metrics"]
		}`),
	}

	result, err := Analyze("Senior Backend Engineer, Go required", testProfile(), caller)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MatchScore != 72 {
		t.Errorf("expected match score 72, got %d", result.MatchScore)
	}
	if len(result.Strengths) != 1 || len(result.Gaps) != 1 || len(result.CoverLetterFocus) != 1 {
		t.Errorf("unexpected result shape: %+v", result)
	}
	if !strings.Contains(caller.capturedSystem, "Jane Doe") {
		t.Error("expected system prompt to embed candidate profile")
	}
	if caller.capturedUser != "Senior Backend Engineer, Go required" {
		t.Errorf("expected job text to be passed as user message, got %q", caller.capturedUser)
	}
}

func TestAnalyze_EmptyJobText(t *testing.T) {
	caller := &fakeCaller{}
	_, err := Analyze("", testProfile(), caller)
	if err == nil {
		t.Fatal("expected error for empty job text")
	}
}

func TestAnalyze_APIError(t *testing.T) {
	caller := &fakeCaller{err: errors.New("network exploded")}
	_, err := Analyze("job text", testProfile(), caller)
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestAnalyze_MalformedToolOutput(t *testing.T) {
	caller := &fakeCaller{response: json.RawMessage(`not json`)}
	_, err := Analyze("job text", testProfile(), caller)
	if err == nil {
		t.Fatal("expected error for malformed tool output")
	}
}

func TestAnalyze_OutOfRangeScore(t *testing.T) {
	caller := &fakeCaller{response: json.RawMessage(`{"match_score": 150, "strengths": [], "gaps": [], "cover_letter_focus": []}`)}
	_, err := Analyze("job text", testProfile(), caller)
	if err == nil {
		t.Fatal("expected error for out-of-range match score")
	}
}
