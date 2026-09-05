## 1. Project setup

- [x] 1.1 Initialize Go module (`go mod init github.com/jdgomez/job-assistant`)
- [x] 1.2 Add `gopkg.in/yaml.v3` dependency
- [x] 1.3 Create package skeleton: `internal/anthropic/`, `internal/profile/`, `internal/analyzer/`, `cmd/`

## 2. Anthropic API client

- [x] 2.1 Define request/response types for the Messages API, including tool-use request (tool schema for `submit_fit_analysis`) and tool-use response parsing
- [x] 2.2 Implement `internal/anthropic.Client` with configurable base URL / HTTP client (for testability)
- [x] 2.3 Implement `Client.CallTool(system, userMessage string, tool ToolSpec) (json.RawMessage, error)` — builds the tool-use request, forces `tool_choice`, parses the tool call output (generic, not job/profile-specific — see note below)
- [x] 2.4 Handle API errors (non-2xx status, network failure, malformed/missing tool call) as returned Go errors, not partial results
- [x] 2.5 Write client tests against an `httptest.Server` returning canned tool-use responses (success, malformed response, HTTP error)

## 3. Candidate profile loading

- [x] 3.1 Define the profile struct (experience, stack, quantified achievements, role preferences) matching `config/profile.yaml`'s shape
- [x] 3.2 Implement `internal/profile.Load(path string) (Profile, error)` using `yaml.v3`, with clear errors for missing file / invalid YAML
- [x] 3.3 Write `config/profile.example.yaml` as a template (no real personal data); document that the real profile file is local/gitignored
- [x] 3.4 Write profile loader tests (valid file, missing file, invalid YAML)

## 4. Analysis logic

- [x] 4.1 Write the system prompt template embedding the candidate profile and analysis instructions (per design.md's tool-use decision)
- [x] 4.2 Implement `internal/analyzer.Analyze(jobText string, profile Profile, client toolCaller) (FitAnalysis, error)` wiring prompt + client together (client is a small interface `*anthropic.Client` satisfies, for testability)
- [x] 4.3 Write analyzer tests using fixture profiles and fixture job postings against a mocked client

## 5. CLI

- [x] 5.1 Implement job posting input resolution: `--job-file` flag if given, else read stdin; error if neither yields content
- [x] 5.2 Wire profile file path (flag or default `config/profile.yaml`) into `internal/profile.Load`
- [x] 5.3 Check `ANTHROPIC_API_KEY` env var before any API call; exit with clear error if unset
- [x] 5.4 Implement plain-text output renderer matching the format in the project notes (match score, strengths, gaps, cover letter focus)
- [x] 5.5 Implement `--json` output renderer
- [x] 5.6 Wire `main.go` / `cmd/analyze.go` together with flag parsing and exit codes

## 6. Docs and examples

- [x] 6.1 Add `examples/sample-job.txt`
- [x] 6.2 Write `README.md`: what it does, requirements (Go version, `ANTHROPIC_API_KEY`), install/usage, example output, known limitations

## 7. Verification

- [x] 7.1 Run `go build` and `go vet` cleanly
- [x] 7.2 Run full test suite (`go test ./...`)
- [x] 7.3 Manual end-to-end run against a real job posting and profile, both text and `--json` output, to confirm output matches the `job-fit-analysis` and `cli-interface` specs
