# job-assistant

A Go CLI that analyzes a job posting against your professional profile using
the Anthropic API, and returns a structured fit analysis: a match score,
strengths to highlight, gaps to address, and recommended cover letter focus
points.

## Why this exists

Job searching means re-reading every posting against your own resume to
judge fit and decide what to emphasize in a cover letter. This tool
automates that first pass so you can spend your time on the applications
worth pursuing.

## Requirements

- Go 1.21+
- An Anthropic API key ([console.anthropic.com](https://console.anthropic.com))

## Installation

```bash
git clone https://github.com/jdgomez/job-assistant.git
cd job-assistant
go build -o job-assistant .
```

## Setup

1. Set your API key:

   ```bash
   export ANTHROPIC_API_KEY=your_key_here
   ```

2. Create your profile from the template:

   ```bash
   cp config/profile.example.yaml config/profile.yaml
   ```

   Edit `config/profile.yaml` with your own experience, skills, quantified
   achievements, and role preferences. This file is gitignored — it's meant
   to hold your personal data locally, not to be committed.

## Usage

```bash
# From a file
./job-assistant --job-file offer.txt

# From a file, JSON output
./job-assistant --job-file offer.txt --json

# From stdin
cat offer.txt | ./job-assistant

# Try it with the bundled example
./job-assistant --job-file examples/sample-job.txt
```

## Example output

```
=== Job Fit Analysis ===

Match score: 72%

Strengths to highlight:
- 13+ years backend experience matches "senior" requirement
- CI/CD pipeline experience (4h → 10min) directly relevant
- Remote async team experience aligns with culture

Gaps detected:
- Go: personal projects only, no commercial experience
- Elixir: not mentioned in profile

Cover letter focus:
- Lead with CI/CD transformation metrics
- Emphasize microservices migration experience
- Address Go gap proactively
```

## Known limitations

- Single-provider (Anthropic only); no support for other LLM providers.
- Each run makes a live API call — there is no caching of repeated analyses.
- No batch mode: one job posting per run.
- No automated cover-letter drafting yet; the tool only recommends focus
  points for a cover letter you write yourself.
- The quality of the analysis depends on how detailed `config/profile.yaml`
  is — a thin profile produces a thin analysis.
