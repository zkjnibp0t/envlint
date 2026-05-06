package profiler_test

import (
	"testing"

	"github.com/user/envlint/profiler"
	"github.com/user/envlint/schema"
)

func makeSchema(vars []schema.VarDef) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestAnalyze_BasicCounts(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	p := profiler.Analyze(env, nil)
	if p.TotalVars != 2 {
		t.Errorf("expected 2 total vars, got %d", p.TotalVars)
	}
	if p.DefinedInSchema != 0 {
		t.Errorf("expected 0 defined in schema, got %d", p.DefinedInSchema)
	}
}

func TestAnalyze_SensitiveVars(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
	}
	p := profiler.Analyze(env, nil)
	if p.SensitiveVars < 2 {
		t.Errorf("expected at least 2 sensitive vars, got %d", p.SensitiveVars)
	}
}

func TestAnalyze_TypeBreakdown(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "PORT", Type: "int"},
		{Name: "DEBUG", Type: "bool"},
	})
	env := map[string]string{
		"PORT":    "8080",
		"DEBUG":   "true",
		"UNKNOWN": "val",
	}
	p := profiler.Analyze(env, s)
	if p.TypeBreakdown["int"] != 1 {
		t.Errorf("expected 1 int, got %d", p.TypeBreakdown["int"])
	}
	if p.TypeBreakdown["bool"] != 1 {
		t.Errorf("expected 1 bool, got %d", p.TypeBreakdown["bool"])
	}
	if p.TypeBreakdown["unknown"] != 1 {
		t.Errorf("expected 1 unknown, got %d", p.TypeBreakdown["unknown"])
	}
}

func TestAnalyze_DefinedVsUndefined(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "APP_NAME", Type: "string"},
	})
	env := map[string]string{
		"APP_NAME": "myapp",
		"EXTRA":    "value",
	}
	p := profiler.Analyze(env, s)
	if p.DefinedInSchema != 1 {
		t.Errorf("expected 1 defined, got %d", p.DefinedInSchema)
	}
	if p.UndefinedInSchema != 1 {
		t.Errorf("expected 1 undefined, got %d", p.UndefinedInSchema)
	}
}

func TestSummary_Format(t *testing.T) {
	p := profiler.Profile{
		TotalVars:         4,
		SensitiveVars:     1,
		DefinedInSchema:   3,
		UndefinedInSchema: 1,
	}
	got := profiler.Summary(p)
	expected := "Total: 4 | Sensitive: 1 | In schema: 3 | Undefined: 1"
	if got != expected {
		t.Errorf("unexpected summary:\ngot:  %s\nwant: %s", got, expected)
	}
}
