package validator_test

import (
	"testing"

	"github.com/user/envlint/schema"
	"github.com/user/envlint/validator"
)

func makeSchema(vars []schema.VarDef) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestValidate_AllPresent(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "APP_NAME", Required: true, Type: "string"},
		{Name: "PORT", Required: true, Type: "int"},
	})
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	report := validator.Validate(s, env)
	if report.HasErrors() {
		t.Errorf("expected no errors, got: %+v", report.Results)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "DATABASE_URL", Required: true, Type: "url"},
	})
	env := map[string]string{}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Error("expected errors for missing required var")
	}
}

func TestValidate_InvalidURL(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "API_URL", Required: true, Type: "url"},
	})
	env := map[string]string{"API_URL": "not-a-url"}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Error("expected URL validation error")
	}
}

func TestValidate_InvalidInt(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "PORT", Required: true, Type: "int"},
	})
	env := map[string]string{"PORT": "abc"}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Error("expected int validation error")
	}
}

func TestValidate_InvalidBool(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "DEBUG", Required: false, Type: "bool"},
	})
	env := map[string]string{"DEBUG": "yes"}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Error("expected bool validation error")
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "API_KEY", Required: true, Type: "string", Pattern: `^[A-Z0-9]{16}$`},
	})
	env := map[string]string{"API_KEY": "ABCD1234EFGH5678"}
	report := validator.Validate(s, env)
	if report.HasErrors() {
		t.Errorf("expected no errors, got: %+v", report.Results)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := makeSchema([]schema.VarDef{
		{Name: "API_KEY", Required: true, Type: "string", Pattern: `^[A-Z0-9]{16}$`},
	})
	env := map[string]string{"API_KEY": "short"}
	report := validator.Validate(s, env)
	if !report.HasErrors() {
		t.Error("expected pattern mismatch error")
	}
}
