## Why

Job searching means manually re-reading each job posting against your own resume to judge fit and figure out what to emphasize in a cover letter. This is repetitive and error-prone, and it's also a concrete, public demonstration of production LLM integration in Go — directly relevant to roles that value hands-on AI tooling experience.

## What Changes

- New Go CLI (`job-assistant`) that takes a job posting and a candidate profile, and returns a structured job-fit analysis.
- Candidate profile is loaded from a YAML config file (`config/profile.yaml`).
- Job posting text can be provided via `--job-file <path>` or piped through stdin.
- Analysis is produced by calling the Anthropic Messages API (`claude-sonnet-4-6`) with a system prompt containing the full candidate profile.
- Output includes: match score (0-100%), top strengths to highlight, gaps to address, and recommended cover letter focus points.
- Output defaults to human-readable plain text; `--json` flag switches to structured JSON.
- The Anthropic API key is read from the `ANTHROPIC_API_KEY` environment variable; missing key is a startup error, not a silent failure.

## Capabilities

### New Capabilities
- `job-fit-analysis`: Given a job posting and a candidate profile, produce a structured fit analysis (score, strengths, gaps, cover letter focus) via the Anthropic API.
- `cli-interface`: Command-line surface for the tool — accepting job posting input (stdin or `--job-file`), loading the candidate profile from config, and rendering output as plain text or JSON (`--json`).

### Modified Capabilities
(none — greenfield project)

## Impact

- New Go module (`github.com/jdgomez/job-assistant`) with a `main.go` entrypoint.
- New dependency: `gopkg.in/yaml.v3` for profile parsing.
- New external dependency: Anthropic Messages API (network call, requires `ANTHROPIC_API_KEY`).
- New config file convention: `config/profile.yaml` (candidate profile, not committed with real personal data — an example/template should be provided).
