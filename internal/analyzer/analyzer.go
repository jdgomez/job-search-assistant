// Package analyzer builds the job-fit prompt and turns the Anthropic API's
// tool-use response into a structured FitAnalysis.
package analyzer

import (
	"encoding/json"
	"fmt"

	"github.com/jdgomez/job-assistant/internal/anthropic"
	"github.com/jdgomez/job-assistant/internal/profile"
)

// FitAnalysis is the structured result of analyzing a job posting against
// a candidate profile.
type FitAnalysis struct {
	MatchScore       int      `json:"match_score"`
	Strengths        []string `json:"strengths"`
	Gaps             []string `json:"gaps"`
	CoverLetterFocus []string `json:"cover_letter_focus"`
}

const toolName = "submit_fit_analysis"

var analysisTool = anthropic.ToolSpec{
	Name:        toolName,
	Description: "Submit the structured job fit analysis for the candidate against the job posting.",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"match_score": map[string]any{
				"type":        "integer",
				"minimum":     0,
				"maximum":     100,
				"description": "Overall fit score from 0 to 100.",
			},
			"strengths": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "Concrete strengths from the candidate's profile that match the job posting's requirements.",
			},
			"gaps": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "Requirements stated in the job posting that the candidate's profile does not address.",
			},
			"cover_letter_focus": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "Recommended points to lead with in a cover letter for this job posting.",
			},
		},
		"required": []string{"match_score", "strengths", "gaps", "cover_letter_focus"},
	},
}

const systemPromptTemplate = `You are a career advisor analyzing job fit for a specific candidate.

Candidate profile:
%s

Analyze the job posting provided by the user and call the %s tool with:
1. A match score (0-100)
2. Strengths to highlight, grounded in the candidate's profile
3. Gaps to address honestly
4. Recommended cover letter focus points

Be specific and actionable. Reference exact requirements from the job posting.`

// toolCaller is the subset of *anthropic.Client that Analyze depends on, so
// tests can substitute a fake instead of calling the real API.
type toolCaller interface {
	CallTool(system, userMessage string, tool anthropic.ToolSpec) (json.RawMessage, error)
}

// Analyze produces a job-fit analysis for jobText against p, using client
// to call the Anthropic API.
func Analyze(jobText string, p profile.Profile, client toolCaller) (FitAnalysis, error) {
	if jobText == "" {
		return FitAnalysis{}, fmt.Errorf("analyzer: job posting text is empty")
	}

	system := fmt.Sprintf(systemPromptTemplate, p.String(), toolName)

	raw, err := client.CallTool(system, jobText, analysisTool)
	if err != nil {
		return FitAnalysis{}, fmt.Errorf("analyzer: %w", err)
	}

	var result FitAnalysis
	if err := json.Unmarshal(raw, &result); err != nil {
		return FitAnalysis{}, fmt.Errorf("analyzer: parsing tool output: %w", err)
	}

	if result.MatchScore < 0 || result.MatchScore > 100 {
		return FitAnalysis{}, fmt.Errorf("analyzer: match score %d out of range [0,100]", result.MatchScore)
	}

	return result, nil
}
