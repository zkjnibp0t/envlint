package profiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlint/envparser"
	"github.com/user/envlint/profiler"
	"github.com/user/envlint/schema"
)

func TestAnalyze_WithExampleFiles(t *testing.T) {
	root := filepath.Join("..", "testdata")

	envPath := filepath.Join(root, "example.env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Skip("testdata/example.env not found, skipping integration test")
	}

	env, err := envparser.Parse(envPath)
	if err != nil {
		t.Fatalf("failed to parse example.env: %v", err)
	}

	schemaPath := filepath.Join(root, "example.schema.yaml")
	var s *schema.Schema
	if _, serr := os.Stat(schemaPath); serr == nil {
		s, err = schema.Load(schemaPath)
		if err != nil {
			t.Fatalf("failed to load schema: %v", err)
		}
	}

	p := profiler.Analyze(env, s)

	if p.TotalVars == 0 {
		t.Error("expected at least one variable in example.env")
	}

	if p.TotalVars != p.DefinedInSchema+p.UndefinedInSchema {
		t.Errorf("defined (%d) + undefined (%d) should equal total (%d)",
			p.DefinedInSchema, p.UndefinedInSchema, p.TotalVars)
	}

	summary := profiler.Summary(p)
	if summary == "" {
		t.Error("expected non-empty summary string")
	}
	t.Logf("Profile summary: %s", summary)
}
