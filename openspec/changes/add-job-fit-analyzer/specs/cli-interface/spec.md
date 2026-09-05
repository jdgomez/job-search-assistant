## Purpose

Provides the command-line surface through which a candidate supplies a job posting and their profile, and receives the fit analysis in a usable format.

## ADDED Requirements

### Requirement: Job posting input
The system SHALL accept the job posting text either from a file path given via a `--job-file` flag, or from standard input when no `--job-file` flag is given.

#### Scenario: Input via file flag
- **WHEN** the tool is invoked with `--job-file offer.txt` and the file exists and is readable
- **THEN** the system uses the file's contents as the job posting text

#### Scenario: Input via stdin
- **WHEN** the tool is invoked without `--job-file` and job posting text is piped to standard input
- **THEN** the system uses the piped content as the job posting text

#### Scenario: No input provided
- **WHEN** the tool is invoked without `--job-file` and standard input is empty or not piped
- **THEN** the system reports an error indicating no job posting was provided, and does not call the analysis API

### Requirement: Candidate profile loading
The system SHALL load the candidate profile from a YAML configuration file before running the analysis, and SHALL report a clear error if the file is missing or cannot be parsed.

#### Scenario: Valid profile file
- **WHEN** the configured profile file exists and is valid YAML matching the expected profile structure
- **THEN** the system loads it and uses it as the candidate profile for analysis

#### Scenario: Missing or invalid profile file
- **WHEN** the configured profile file does not exist, or exists but is not valid YAML
- **THEN** the system reports an error identifying the problem and exits without calling the analysis API

### Requirement: Output format selection
The system SHALL render the analysis as human-readable plain text by default, and as structured JSON when the `--json` flag is given.

#### Scenario: Default text output
- **WHEN** the tool is run without `--json`
- **THEN** the system prints the match score, strengths, gaps, and cover letter focus points as formatted plain text

#### Scenario: JSON output
- **WHEN** the tool is run with `--json`
- **THEN** the system prints a single JSON object containing the match score, strengths, gaps, and cover letter focus points

### Requirement: Missing API key is a startup error
The system SHALL check for the `ANTHROPIC_API_KEY` environment variable before attempting analysis, and SHALL exit with a clear error message if it is not set, rather than failing later with an opaque API error.

#### Scenario: API key not set
- **WHEN** the tool is invoked and the `ANTHROPIC_API_KEY` environment variable is not set
- **THEN** the system exits with a non-zero status and an error message naming the missing environment variable, without making a network request
