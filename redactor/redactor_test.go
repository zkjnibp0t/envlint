package redactor_test

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"envlint/redactor"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_PASSWORD":  "hunter2",
		"API_KEY":      "abc123",
		"APP_PORT":     "8080",
		"SERVICE_NAME": "envlint",
		"JWT_SECRET":   "topsecret",
	}
}

func TestRedact_SensitiveKeysReplaced(t *testing.T) {
	result := redactor.Redact(sampleEnv(), redactor.Options{})

	sensitive := []string{"DB_PASSWORD", "API_KEY", "JWT_SECRET"}
	for _, k := range sensitive {
		if result[k] != "[REDACTED]" {
			t.Errorf("expected %s to be redacted, got %q", k, result[k])
		}
	}
}

func TestRedact_SafeKeysPreserved(t *testing.T) {
	result := redactor.Redact(sampleEnv(), redactor.Options{})

	if result["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT to be preserved, got %q", result["APP_PORT"])
	}
	if result["SERVICE_NAME"] != "envlint" {
		t.Errorf("expected SERVICE_NAME to be preserved, got %q", result["SERVICE_NAME"])
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	result := redactor.Redact(sampleEnv(), redactor.Options{Placeholder: "***"})
	if result["DB_PASSWORD"] != "***" {
		t.Errorf("expected custom placeholder, got %q", result["DB_PASSWORD"])
	}
}

func TestRedact_ExtraRules(t *testing.T) {
	env := map[string]string{"STRIPE_PRIV": "sk_live_xxx", "APP_NAME": "myapp"}
	opts := redactor.Options{
		ExtraRules: []redactor.Rule{{KeyPattern: "priv", Replacement: ""}},
	}
	result := redactor.Redact(env, opts)
	if result["STRIPE_PRIV"] != "[REDACTED]" {
		t.Errorf("expected STRIPE_PRIV to be redacted via extra rule")
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be preserved")
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	original := sampleEnv()
	redactor.Redact(original, redactor.Options{})
	if original["DB_PASSWORD"] != "hunter2" {
		t.Error("original map was mutated")
	}
}

func TestKeys_ReturnsSensitiveKeys(t *testing.T) {
	keys := redactor.Keys(sampleEnv(), redactor.Options{})
	sort.Strings(keys)
	expected := []string{"API_KEY", "DB_PASSWORD", "JWT_SECRET"}
	if strings.Join(keys, ",") != strings.Join(expected, ",") {
		t.Errorf("Keys() = %v, want %v", keys, expected)
	}
}

func TestWriteReport_WithRedacted(t *testing.T) {
	original := sampleEnv()
	redacted := redactor.Redact(original, redactor.Options{})
	var buf bytes.Buffer
	redactor.WriteReport(&buf, original, redacted)
	out := buf.String()
	if !strings.Contains(out, "3 sensitive key(s) redacted") {
		t.Errorf("unexpected report output: %s", out)
	}
}

func TestWriteReport_NoSensitiveKeys(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080"}
	var buf bytes.Buffer
	redactor.WriteReport(&buf, env, env)
	if !strings.Contains(buf.String(), "no sensitive keys found") {
		t.Errorf("expected clean report")
	}
}
