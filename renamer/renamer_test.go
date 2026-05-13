package renamer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envlint/renamer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_PORT":    "8080",
	}
}

func TestRename_BasicRename(t *testing.T) {
	env := baseEnv()
	rules := []renamer.Rule{{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"}}
	out, res := renamer.Rename(env, rules)

	if _, ok := out["DATABASE_HOST"]; !ok {
		t.Error("expected DATABASE_HOST to exist")
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
	if len(res.Renamed) != 1 {
		t.Errorf("expected 1 renamed, got %d", len(res.Renamed))
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	env := baseEnv()
	rules := []renamer.Rule{{OldKey: "NONEXISTENT", NewKey: "SOMETHING"}}
	_, res := renamer.Rename(env, rules)

	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected 0 renamed, got %d", len(res.Renamed))
	}
}

func TestRename_ConflictPreservesOldKey(t *testing.T) {
	env := baseEnv()
	// APP_PORT already exists; renaming DB_HOST → APP_PORT should conflict.
	rules := []renamer.Rule{{OldKey: "DB_HOST", NewKey: "APP_PORT"}}
	out, res := renamer.Rename(env, rules)

	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if out["DB_HOST"] != "localhost" {
		t.Error("expected DB_HOST to be preserved after conflict")
	}
}

func TestRename_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	rules := []renamer.Rule{{OldKey: "APP_PORT", NewKey: "SERVER_PORT"}}
	renamer.Rename(env, rules)

	if _, ok := env["APP_PORT"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := renamer.ParseRules([]string{"OLD_KEY=NEW_KEY", "FOO=BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].OldKey != "OLD_KEY" || rules[0].NewKey != "NEW_KEY" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := renamer.ParseRules([]string{"NOEQUALSSIGN"})
	if err == nil {
		t.Error("expected error for malformed rule")
	}
}

func TestWriteReport_ShowsAllSections(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	rules := []renamer.Rule{
		{OldKey: "A", NewKey: "ALPHA"},    // renamed
		{OldKey: "MISSING", NewKey: "X"}, // skipped
		{OldKey: "B", NewKey: "ALPHA"},   // conflict (ALPHA now exists)
	}
	_, res := renamer.Rename(env, rules)

	var buf bytes.Buffer
	renamer.WriteReport(&buf, res)
	out := buf.String()

	if !strings.Contains(out, "Renamed") {
		t.Error("expected Renamed section")
	}
	if !strings.Contains(out, "Skipped") {
		t.Error("expected Skipped section")
	}
	if !strings.Contains(out, "Summary") {
		t.Error("expected Summary line")
	}
}
