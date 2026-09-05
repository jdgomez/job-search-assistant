package profile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and parses the candidate profile YAML file at path.
func Load(path string) (Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, fmt.Errorf("profile: reading %s: %w", path, err)
	}

	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Profile{}, fmt.Errorf("profile: parsing %s: %w", path, err)
	}

	return p, nil
}
