package normalizer_test

import (
	"testing"

	"github.com/envlint/envlint/normalizer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db_host":  "localhost",
		"DB_PORT":  "5432",
		"api_key":  "  secret  ",
		"App_Name": "envlint",
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	opts := normalizer.Options{UppercaseKeys: true, TrimValues: false, DeduplicateKeys: false}
	res := normalizer.Normalize(baseEnv(), opts)

	for k := range res.Env {
		if k != "DB_HOST" && k != "DB_PORT" && k != "API_KEY" && k != "APP_NAME" {
			t.Errorf("unexpected key in result: %q", k)
		}
	}
	if len(res.Renamed) == 0 {
		t.Error("expected renamed keys, got none")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	opts := normalizer.Options{UppercaseKeys: false, TrimValues: true, DeduplicateKeys: false}
	res := normalizer.Normalize(baseEnv(), opts)

	if got := res.Env["api_key"]; got != "secret" {
		t.Errorf("expected trimmed value %q, got %q", "secret", got)
	}
}

func TestNormalize_NoTrimValues(t *testing.T) {
	opts := normalizer.Options{UppercaseKeys: false, TrimValues: false, DeduplicateKeys: false}
	res := normalizer.Normalize(baseEnv(), opts)

	if got := res.Env["api_key"]; got != "  secret  " {
		t.Errorf("expected untrimmed value, got %q", got)
	}
}

func TestNormalize_DeduplicateKeys(t *testing.T) {
	env := map[string]string{
		"HOST": "first",
	}
	// Simulate a second entry by merging manually
	env2 := map[string]string{
		"HOST": "first",
		"host": "second",
	}
	opts := normalizer.Options{UppercaseKeys: true, TrimValues: false, DeduplicateKeys: true}
	res := normalizer.Normalize(env2, opts)

	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected HOST key in result")
	}
	if len(res.Duplicates) == 0 {
		t.Error("expected duplicate to be recorded")
	}
	_ = env
}

func TestNormalize_DefaultOptions(t *testing.T) {
	opts := normalizer.DefaultOptions()
	res := normalizer.Normalize(baseEnv(), opts)

	if v, ok := res.Env["API_KEY"]; !ok || v != "secret" {
		t.Errorf("expected API_KEY=secret, got %q (present=%v)", v, ok)
	}
	if _, ok := res.Env["db_host"]; ok {
		t.Error("lowercase key should have been uppercased")
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	opts := normalizer.DefaultOptions()
	normalizer.Normalize(env, opts)

	if _, ok := env["db_host"]; !ok {
		t.Error("original map should not be mutated")
	}
}
