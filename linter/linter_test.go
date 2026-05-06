package linter_test

import (
	"testing"

	"github.com/user/envlint/linter"
	"github.com/user/envlint/schema"
)

func makeSchema() *schema.Schema {
	return &schema.Schema{
		Vars: []schema.Var{
			{Name: "PORT", Type: "int", Required: true},
			{Name: "DATABASE_URL", Type: "url", Required: true},
			{Name: "OLD_KEY", Type: "string", Deprecated: true},
		},
	}
}

func TestRun_AllValid(t *testing.T) {
	env := map[string]string{
		"PORT":         "8080",
		"DATABASE_URL": "http://localhost:5432",
	}
	r := linter.Run(env, makeSchema())
	if r.HasErrors {
		t.Errorf("expected no errors, got validation=%v audit=%v", r.ValidationErrors, r.AuditIssues)
	}
	if len(r.ValidationErrors) != 0 {
		t.Errorf("expected 0 validation errors, got %d", len(r.ValidationErrors))
	}
}

func TestRun_MissingRequired(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
	}
	r := linter.Run(env, makeSchema())
	if !r.HasErrors {
		t.Error("expected HasErrors=true for missing required var")
	}
	if len(r.ValidationErrors) == 0 {
		t.Error("expected at least one validation error")
	}
}

func TestRun_InvalidType(t *testing.T) {
	env := map[string]string{
		"PORT":         "not-a-number",
		"DATABASE_URL": "http://localhost",
	}
	r := linter.Run(env, makeSchema())
	if !r.HasErrors {
		t.Error("expected HasErrors=true for invalid int")
	}
}

func TestRun_DeprecatedKeyInUse(t *testing.T) {
	env := map[string]string{
		"PORT":         "3000",
		"DATABASE_URL": "http://db:5432",
		"OLD_KEY":      "some-value",
	}
	r := linter.Run(env, makeSchema())
	if len(r.AuditIssues) == 0 {
		t.Error("expected at least one audit issue for deprecated key")
	}
}

func TestRun_SuggestionsPopulated(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
		// DATABSE_URL is a typo of DATABASE_URL
		"DATABSE_URL": "http://localhost",
	}
	r := linter.Run(env, makeSchema())
	if len(r.Suggestions) == 0 {
		t.Error("expected suggestions for typo'd key")
	}
}
