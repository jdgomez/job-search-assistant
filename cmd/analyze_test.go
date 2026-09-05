package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdgomez/job-assistant/internal/analyzer"
)

func TestResolveJobText_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "job.txt")
	if err := os.WriteFile(path, []byte("Senior Backend Engineer"), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	text, err := resolveJobText(Options{JobFile: path}, strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "Senior Backend Engineer" {
		t.Errorf("expected file contents, got %q", text)
	}
}

func TestResolveJobText_MissingFile(t *testing.T) {
	_, err := resolveJobText(Options{JobFile: "/does/not/exist.txt"}, strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveJobText_FromStdin(t *testing.T) {
	text, err := resolveJobText(Options{}, strings.NewReader("Senior Backend Engineer"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "Senior Backend Engineer" {
		t.Errorf("expected stdin contents, got %q", text)
	}
}

func TestResolveJobText_EmptyStdin(t *testing.T) {
	_, err := resolveJobText(Options{}, strings.NewReader("   \n"))
	if err == nil {
		t.Fatal("expected error for empty stdin")
	}
}

func testResult() analyzer.FitAnalysis {
	return analyzer.FitAnalysis{
		MatchScore:       72,
		Strengths:        []string{"13+ years backend experience"},
		Gaps:             []string{"No commercial Go experience"},
		CoverLetterFocus: []string{"Lead with CI/CD transformation metrics"},
	}
}

func TestRenderText(t *testing.T) {
	var buf bytes.Buffer
	if err := renderText(&buf, testResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Match score: 72%", "13+ years backend experience", "No commercial Go experience", "Lead with CI/CD transformation metrics"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := renderJSON(&buf, testResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded analyzer.FitAnalysis
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded.MatchScore != 72 {
		t.Errorf("expected match score 72, got %d", decoded.MatchScore)
	}
}

func TestRun_MissingAPIKey(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")

	err := Run(Options{}, strings.NewReader("job text"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error when ANTHROPIC_API_KEY is unset")
	}
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Errorf("expected error to mention ANTHROPIC_API_KEY, got: %v", err)
	}
}
