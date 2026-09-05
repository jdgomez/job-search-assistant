## Purpose

Turns a job posting and a candidate profile into a structured, actionable fit assessment, so the candidate can quickly decide whether to apply and how to position their application.

## ADDED Requirements

### Requirement: Structured fit analysis
The system SHALL analyze a given job posting against a given candidate profile and produce: a match score between 0 and 100 (inclusive), a list of strengths to highlight, a list of gaps to address, and a list of recommended cover letter focus points.

#### Scenario: Well-matched candidate
- **WHEN** the candidate profile covers most of the job posting's stated requirements
- **THEN** the system returns a match score, at least one strength, and cover letter focus points referencing the candidate's relevant experience

#### Scenario: Poorly-matched candidate
- **WHEN** the candidate profile is missing several requirements stated in the job posting
- **THEN** the system returns a lower match score and lists the missing requirements as gaps

### Requirement: Analysis grounded in candidate profile
The system SHALL include the complete candidate profile as context when generating the analysis, so results are specific to that candidate rather than generic.

#### Scenario: Profile content reflected in output
- **WHEN** the candidate profile lists a specific quantified achievement (e.g. a metric or outcome)
- **THEN** the analysis may reference that achievement in the strengths or cover letter focus points when relevant to the job posting

### Requirement: Requirements-referenced output
The system SHALL produce gaps and strengths that reference concrete requirements stated in the job posting, rather than generic statements unconnected to the posting's content.

#### Scenario: Gap tied to explicit requirement
- **WHEN** the job posting states a requirement not present anywhere in the candidate profile
- **THEN** the system lists that specific requirement as a gap

### Requirement: Analysis failure is surfaced, not silenced
The system SHALL report an error to the caller when the fit analysis cannot be produced (e.g. the Anthropic API call fails, times out, or returns an unusable response), and SHALL NOT return a fabricated or partial analysis as if it were complete.

#### Scenario: API call fails
- **WHEN** the Anthropic API request fails (network error, non-success HTTP status, or malformed response)
- **THEN** the system returns an error describing the failure and produces no analysis output
