// Command job-assistant analyzes a job posting against a candidate profile
// using the Anthropic API and prints a structured fit analysis.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jdgomez/job-assistant/cmd"
)

func main() {
	jobFile := flag.String("job-file", "", "Path to a file containing the job posting text. If omitted, reads from stdin.")
	profilePath := flag.String("profile", "config/profile.yaml", "Path to the candidate profile YAML file.")
	jsonOutput := flag.Bool("json", false, "Output the analysis as JSON instead of plain text.")
	flag.Parse()

	opts := cmd.Options{
		JobFile:     *jobFile,
		ProfilePath: *profilePath,
		JSON:        *jsonOutput,
	}

	if err := cmd.Run(opts, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
