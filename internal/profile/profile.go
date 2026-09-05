// Package profile loads and represents the candidate profile used as
// context for job-fit analysis.
package profile

import "strings"

// Profile is a candidate's professional profile.
type Profile struct {
	Name         string       `yaml:"name"`
	Headline     string       `yaml:"headline"`
	Summary      string       `yaml:"summary"`
	Experience   []Experience `yaml:"experience"`
	Skills       []string     `yaml:"skills"`
	Achievements []string     `yaml:"achievements"`
	Preferences  Preferences  `yaml:"preferences"`
}

// Experience is a single role in the candidate's work history.
type Experience struct {
	Role        string `yaml:"role"`
	Company     string `yaml:"company"`
	Years       string `yaml:"years"`
	Description string `yaml:"description"`
}

// Preferences captures the kind of role the candidate is looking for.
type Preferences struct {
	Roles     []string `yaml:"roles"`
	Remote    bool     `yaml:"remote"`
	Locations []string `yaml:"locations"`
}

// String renders the profile as plain text suitable for embedding in an
// LLM prompt.
func (p Profile) String() string {
	var b strings.Builder

	if p.Name != "" {
		b.WriteString("Name: " + p.Name + "\n")
	}
	if p.Headline != "" {
		b.WriteString("Headline: " + p.Headline + "\n")
	}
	if p.Summary != "" {
		b.WriteString("Summary: " + p.Summary + "\n")
	}

	if len(p.Experience) > 0 {
		b.WriteString("\nExperience:\n")
		for _, e := range p.Experience {
			b.WriteString(fmtExperience(e))
		}
	}

	if len(p.Skills) > 0 {
		b.WriteString("\nSkills: " + strings.Join(p.Skills, ", ") + "\n")
	}

	if len(p.Achievements) > 0 {
		b.WriteString("\nQuantified achievements:\n")
		for _, a := range p.Achievements {
			b.WriteString("- " + a + "\n")
		}
	}

	if len(p.Preferences.Roles) > 0 || len(p.Preferences.Locations) > 0 {
		b.WriteString("\nRole preferences:\n")
		if len(p.Preferences.Roles) > 0 {
			b.WriteString("- Roles: " + strings.Join(p.Preferences.Roles, ", ") + "\n")
		}
		b.WriteString("- Remote: " + boolToYesNo(p.Preferences.Remote) + "\n")
		if len(p.Preferences.Locations) > 0 {
			b.WriteString("- Locations: " + strings.Join(p.Preferences.Locations, ", ") + "\n")
		}
	}

	return b.String()
}

func fmtExperience(e Experience) string {
	var b strings.Builder
	b.WriteString("- " + e.Role)
	if e.Company != "" {
		b.WriteString(" at " + e.Company)
	}
	if e.Years != "" {
		b.WriteString(" (" + e.Years + ")")
	}
	b.WriteString("\n")
	if e.Description != "" {
		b.WriteString("  " + e.Description + "\n")
	}
	return b.String()
}

func boolToYesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
