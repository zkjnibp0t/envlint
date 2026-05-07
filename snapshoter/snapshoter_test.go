package snapshoter_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envlint/snapshoter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"API_KEY":  "secret",
	}
}

func TestCapture_CopiesEnv(t *testing.T) {
	env := sampleEnv()
	snap := snapshoter.Capture(env, "test")

	if snap.Label != "test" {
		t.Errorf("expected label 'test', got %q", snap.Label)
	}
	if snap.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production")
	}
	// Mutating original should not affect snapshot
	env["APP_ENV"] = "staging"
	if snap.Env["APP_ENV"] != "production" {
		t.Error("snapshot was mutated by original map change")
	}
}

func TestCapture_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	snap := snapshoter.Capture(sampleEnv(), "ts-test")
	after := time.Now().UTC()

	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("captured_at %v not in expected range", snap.CapturedAt)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshoter.Capture(sampleEnv(), "round-trip")
	if err := snapshoter.Save(orig, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshoter.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Label != orig.Label {
		t.Errorf("label mismatch: got %q, want %q", loaded.Label, orig.Label)
	}
	for k, v := range orig.Env {
		if loaded.Env[k] != v {
			t.Errorf("key %s: got %q, want %q", k, loaded.Env[k], v)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshoter.Load("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)

	_, err := snapshoter.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestKeys_Sorted(t *testing.T) {
	snap := snapshoter.Capture(sampleEnv(), "keys")
	keys := snap.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}
