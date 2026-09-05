## Context

Greenfield Go CLI, no existing code. See proposal.md - Why/What Changes for motivation and scope. The only non-trivial technical question is how to get a reliably structured response out of the Anthropic API, since the CLI needs to parse match score / strengths / gaps / cover letter focus into both plain-text and JSON output.

## Goals / Non-Goals

**Goals:**
- Reliable, parseable structured output from a single Anthropic API call per analysis.
- Clear separation between API client, profile loading, and analysis/prompt logic (per proposal's package layout).
- Testable without hitting the real Anthropic API.

**Non-Goals:**
- Caching, batch processing, or cover-letter drafting (listed as post-v1 extensions in the original project notes; not part of this change).
- Multi-provider LLM support — Anthropic only.

## Decisions

### Structured output via forced tool use, not prompted JSON
Use the Anthropic Messages API's tool-use feature: define a single tool (e.g. `submit_fit_analysis`) with a JSON schema for `{match_score, strengths[], gaps[], cover_letter_focus[]}`, and force its use via `tool_choice`. The tool call's `input` is the structured result — no free-text parsing needed.

Alternative considered: ask the model to "return JSON" in the prompt and parse the text response. Rejected — free-text JSON is prone to prose wrapping, markdown code fences, or truncation, which would require fragile cleanup/parsing and produce occasional runtime failures that are hard to test for.

### stdlib HTTP client, no Anthropic SDK dependency
`internal/anthropic/client.go` builds the request and parses the response using `net/http` and `encoding/json` directly, per the original plan. Keeps dependencies minimal (only `gopkg.in/yaml.v3` beyond stdlib) and keeps the tool-use request/response shape fully under our control for testing.

Alternative considered: use an official/community Anthropic Go SDK. Rejected for v1 — the API surface needed (one endpoint, one tool-use call) is small enough that a hand-rolled client is simpler to reason about and test than adopting a new dependency.

### Package boundaries
Follow the structure in the original project notes:
- `internal/anthropic`: HTTP client + request/response types for the Messages API. No knowledge of profiles or CLI flags.
- `internal/profile`: loads and validates `config/profile.yaml` into a Go struct. No knowledge of the API.
- `internal/analyzer`: owns the system prompt template and the tool schema; takes job text + profile struct, calls `internal/anthropic`, returns a parsed `FitAnalysis` struct.
- `cmd/analyze.go` + `main.go`: flag parsing, stdin handling, output rendering (text vs JSON).

This keeps `internal/anthropic` and `internal/profile` mockable/testable in isolation from prompt logic and CLI concerns.

### Testing the API-dependent path
`internal/anthropic.Client` takes an `http.Client`-like interface (or base URL) so tests can point it at an `httptest.Server` returning canned tool-use responses, instead of calling the real API. Analyzer-level tests exercise prompt construction and response-to-`FitAnalysis` mapping against fixture responses, not live calls.

## Risks / Trade-offs

- **[Risk]** Forced tool-use may make the model's phrasing feel slightly more terse/templated than free-form prose. → Mitigation: the tool's parameter descriptions instruct the model to be specific and reference exact job posting requirements (per the `job-fit-analysis` spec's requirements-referenced-output requirement).
- **[Risk]** Every analysis run costs an API call (money + latency); no caching in v1. → Mitigation: explicitly out of scope for this change (see Non-Goals); acceptable for a single-candidate CLI tool used a handful of times per application.
- **[Risk]** `config/profile.yaml` may contain real personal data if a real profile is committed by accident. → Mitigation: ship only an example/template profile in the repo; document the real profile as a local, gitignored file.

## Migration Plan

None — new tool, no existing users or data to migrate.
