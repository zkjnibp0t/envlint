package masker_test

import (
	"testing"

	"github.com/yourorg/envlint/masker"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"GITHUB_TOKEN",
		"AWS_SECRET",
		"PRIVATE_KEY",
		"AUTH_TOKEN",
		"app_secret",    // lowercase should still match
		"my_credentials",
	}
	for _, key := range sensitiveKeys {
		if !masker.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_IgnoresSafeKeys(t *testing.T) {
	safeKeys := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_URL",
		"REDIS_HOST",
	}
	for _, key := range safeKeys {
		if masker.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMaskEnv_RedactsSensitiveValues(t *testing.T) {
	input := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}
	masked := masker.MaskEnv(input)

	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked, got %q", masked["APP_ENV"])
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should not be masked, got %q", masked["PORT"])
	}
	if masked["DB_PASSWORD"] != masker.Masked {
		t.Errorf("DB_PASSWORD should be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["API_KEY"] != masker.Masked {
		t.Errorf("API_KEY should be masked, got %q", masked["API_KEY"])
	}
}

func TestMaskEnv_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
	}
	_ = masker.MaskEnv(input)
	if input["DB_PASSWORD"] != "supersecret" {
		t.Error("MaskEnv must not mutate the original map")
	}
}

func TestMaskValue_Sensitive(t *testing.T) {
	got := masker.MaskValue("GITHUB_TOKEN", "ghp_abc123")
	if got != masker.Masked {
		t.Errorf("expected %q, got %q", masker.Masked, got)
	}
}

func TestMaskValue_NotSensitive(t *testing.T) {
	got := masker.MaskValue("APP_ENV", "staging")
	if got != "staging" {
		t.Errorf("expected %q, got %q", "staging", got)
	}
}
