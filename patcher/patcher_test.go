package patcher_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "development",
		"LOG_LEVEL": "info",
		"PORT":     "8080",
	}
}

func TestPatch_AddsNewKey(t *testing.T) {
	res := patcher.Patch(baseEnv(), map[string]string{"NEW_KEY": "hello"})
	if res.Env["NEW_KEY"] != "hello" {
		t.Fatalf("expected NEW_KEY=hello, got %q", res.Env["NEW_KEY"])
	}
	found := false
	for _, c := range res.Changes {
		if c.Key == "NEW_KEY" && c.Op == patcher.OpAdded {
			found = true
		}
	}
	if !found {
		t.Error("expected OpAdded change for NEW_KEY")
	}
}

func TestPatch_UpdatesExistingKey(t *testing.T) {
	res := patcher.Patch(baseEnv(), map[string]string{"PORT": "9090"})
	if res.Env["PORT"] != "9090" {
		t.Fatalf("expected PORT=9090, got %q", res.Env["PORT"])
	}
	for _, c := range res.Changes {
		if c.Key == "PORT" {
			if c.Op != patcher.OpUpdated {
				t.Errorf("expected OpUpdated, got %q", c.Op)
			}
			if c.OldValue != "8080" || c.NewValue != "9090" {
				t.Errorf("unexpected old/new values: %q -> %q", c.OldValue, c.NewValue)
			}
		}
	}
}

func TestPatch_UnchangedWhenValueSame(t *testing.T) {
	res := patcher.Patch(baseEnv(), map[string]string{"APP_ENV": "development"})
	for _, c := range res.Changes {
		if c.Key == "APP_ENV" && c.Op != patcher.OpUnchanged {
			t.Errorf("expected OpUnchanged for APP_ENV, got %q", c.Op)
		}
	}
}

func TestPatch_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	patcher.Patch(env, map[string]string{"PORT": "1234", "EXTRA": "yes"})
	if env["PORT"] != "8080" {
		t.Error("original env was mutated")
	}
	if _, ok := env["EXTRA"]; ok {
		t.Error("original env gained EXTRA key")
	}
}

func TestWriteReport_NoChanges(t *testing.T) {
	env := baseEnv()
	res := patcher.Patch(env, map[string]string{"APP_ENV": "development"})
	var buf bytes.Buffer
	patcher.WriteReport(&buf, res)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' message, got: %s", buf.String())
	}
}

func TestWriteReport_WithChanges(t *testing.T) {
	res := patcher.Patch(baseEnv(), map[string]string{"PORT": "9090", "BRAND_NEW": "val"})
	var buf bytes.Buffer
	patcher.WriteReport(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "1 added") {
		t.Errorf("expected '1 added' in output: %s", out)
	}
	if !strings.Contains(out, "1 updated") {
		t.Errorf("expected '1 updated' in output: %s", out)
	}
	if !strings.Contains(out, "BRAND_NEW") {
		t.Errorf("expected BRAND_NEW in output: %s", out)
	}
}
