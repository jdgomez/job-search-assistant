package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.yaml")
	content := `
name: Jane Doe
headline: Senior Backend Engineer
skills:
  - Go
  - PostgreSQL
achievements:
  - "Cut CI pipeline time from 4h to 10min"
preferences:
  roles:
    - Backend Engineer
  remote: true
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	p, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Jane Doe" {
		t.Errorf("expected name %q, got %q", "Jane Doe", p.Name)
	}
	if len(p.Skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(p.Skills))
	}
	if !p.Preferences.Remote {
		t.Error("expected Preferences.Remote to be true")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.yaml")
	if err := os.WriteFile(path, []byte("name: [unclosed"), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
