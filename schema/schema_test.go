package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlint/schema"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return p
}

func TestLoad_ValidSchema(t *testing.T) {
	raw := `
vars:
  DATABASE_URL:
    required: true
    type: url
    description: Primary database connection string
  DEBUG:
    required: false
    type: bool
    default: "false"
`
	p := writeTemp(t, raw)
	s, err := schema.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}

	dbURL := s.Vars["DATABASE_URL"]
	if !dbURL.Required {
		t.Error("DATABASE_URL should be required")
	}
	if dbURL.Type != schema.TypeURL {
		t.Errorf("expected type url, got %s", dbURL.Type)
	}
}

func TestLoad_EmptyVars(t *testing.T) {
	p := writeTemp(t, "vars:\n")
	s, err := schema.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Vars == nil {
		t.Error("Vars map should not be nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := schema.Load("/nonexistent/path/schema.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTemp(t, ": invalid: yaml: [")
	_, err := schema.Load(p)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}
