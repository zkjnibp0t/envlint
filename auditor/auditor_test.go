package auditor_test

import (
	"testing"

	"github.com/user/envlint/auditor"
	"github.com/user/envlint/schema"
)

func makeSchema(vars []schema.VarDefinition) *schema.Schema {
	return &schema.Schema{Vars: vars}
}

func TestAudit_UnusedKey(t *testing.T) {
	s := makeSchema([]schema.VarDefinition{
		{Name: "APP_ENV", Type: "string"},
	})
	env := map[string]string{
		"APP_ENV":  "production",
		"LEFTOVER": "value",
	}

	result := auditor.Audit(env, s)

	if len(result.Unused) != 1 {
		t.Fatalf("expected 1 unused issue, got %d", len(result.Unused))
	}
	if result.Unused[0].Key != "LEFTOVER" {
		t.Errorf("expected LEFTOVER, got %s", result.Unused[0].Key)
	}
}

func TestAudit_DeprecatedInUse(t *testing.T) {
	s := makeSchema([]schema.VarDefinition{
		{Name: "OLD_API_KEY", Type: "string", Deprecated: true, DeprecationNote: "use NEW_API_KEY instead"},
	})
	env := map[string]string{
		"OLD_API_KEY": "abc123",
	}

	result := auditor.Audit(env, s)

	if len(result.Deprecated) != 1 {
		t.Fatalf("expected 1 deprecated issue, got %d", len(result.Deprecated))
	}
	if result.Deprecated[0].Key != "OLD_API_KEY" {
		t.Errorf("expected OLD_API_KEY, got %s", result.Deprecated[0].Key)
	}
}

func TestAudit_DeprecatedNotInUse(t *testing.T) {
	s := makeSchema([]schema.VarDefinition{
		{Name: "OLD_KEY", Type: "string", Deprecated: true},
	})
	env := map[string]string{}

	result := auditor.Audit(env, s)

	if len(result.Deprecated) != 0 {
		t.Errorf("expected no deprecated issues, got %d", len(result.Deprecated))
	}
}

func TestAudit_NoIssues(t *testing.T) {
	s := makeSchema([]schema.VarDefinition{
		{Name: "APP_PORT", Type: "int"},
		{Name: "APP_ENV", Type: "string"},
	})
	env := map[string]string{
		"APP_PORT": "8080",
		"APP_ENV":  "staging",
	}

	result := auditor.Audit(env, s)

	if result.HasIssues() {
		t.Errorf("expected no issues, got unused=%d deprecated=%d",
			len(result.Unused), len(result.Deprecated))
	}
}

func TestAudit_HasIssues(t *testing.T) {
	s := makeSchema([]schema.VarDefinition{})
	env := map[string]string{"GHOST": "value"}

	result := auditor.Audit(env, s)

	if !result.HasIssues() {
		t.Error("expected HasIssues to return true")
	}
}
