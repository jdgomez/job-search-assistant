// Package cmd implements the job-assistant CLI's analyze command.
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jdgomez/job-assistant/internal/analyzer"
	"github.com/jdgomez/job-assistant/internal/anthropic"
	"github.com/jdgomez/job-assistant/internal/profile"
)

// Options holds parsed CLI flags for the analyze command.
type Options struct {
	JobFile     string
	ProfilePath string
	JSON        bool
}

// Run executes the analyze command: it checks for an API key, resolves the
// job posting input, loads the candidate profile, calls the Anthropic API,
// and writes the resulting analysis to stdout.
func Run(opts Options, stdin io.Reader, stdout io.Writer) error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	jobText, err := resolveJobText(opts, stdin)
	if err != nil {
		return err
	}

	p, err := profile.Load(opts.ProfilePath)
	if err != nil {
		return err
	}

	client := anthropic.NewClient(apiKey)
	result, err := analyzer.Analyze(jobText, p, client)
	if err != nil {
		return err
	}

	if opts.JSON {
		return renderJSON(stdout, result)
	}
	return renderText(stdout, result)
}

func resolveJobText(opts Options, stdin io.Reader) (string, error) {
	if opts.JobFile != "" {
		data, err := os.ReadFile(opts.JobFile)
		if err != nil {
			return "", fmt.Errorf("reading job file %s: %w", opts.JobFile, err)
		}
		return string(data), nil
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("reading job posting from stdin: %w", err)
	}
	text := string(data)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("no job posting provided: pass --job-file or pipe text via stdin")
	}
	return text, nil
}

func renderJSON(w io.Writer, result analyzer.FitAnalysis) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func renderText(w io.Writer, result analyzer.FitAnalysis) error {
	fmt.Fprintln(w, "=== Job Fit Analysis ===")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Match score: %d%%\n", result.MatchScore)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Strengths to highlight:")
	for _, s := range result.Strengths {
		fmt.Fprintf(w, "- %s\n", s)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Gaps detected:")
	for _, g := range result.Gaps {
		fmt.Fprintf(w, "- %s\n", g)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Cover letter focus:")
	for _, c := range result.CoverLetterFocus {
		fmt.Fprintf(w, "- %s\n", c)
	}

	return nil
}
