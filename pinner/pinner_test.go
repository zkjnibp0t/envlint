package pinner

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"APP_PORT": "8080",
		"DB_HOST":  "db.example.com",
	}
}

func TestPinAll_ReturnsAllKeys(t *testing.T) {
	env := baseEnv()
	pins := PinAll(env)
	if len(pins) != len(env) {
		t.Fatalf("expected %d pins, got %d", len(env), len(pins))
	}
	for _, p := range pins {
		if p.PinnedAt.IsZero() {
			t.Errorf("PinnedAt not set for key %s", p.Key)
		}
		if p.PinnedAt.After(time.Now().UTC().Add(time.Second)) {
			t.Errorf("PinnedAt is in the future for key %s", p.Key)
		}
	}
}

func TestCheck_AllMatch(t *testing.T) {
	env := baseEnv()
	pins := PinAll(env)
	r := Check(env, pins)
	if len(r.Drifted) != 0 {
		t.Errorf("expected no drift, got %d drifted", len(r.Drifted))
	}
	if len(r.Matched) != len(env) {
		t.Errorf("expected %d matched, got %d", len(env), len(r.Matched))
	}
}

func TestCheck_DetectsDrift(t *testing.T) {
	env := baseEnv()
	pins := PinAll(env)
	env["APP_PORT"] = "9090"
	r := Check(env, pins)
	if len(r.Drifted) != 1 {
		t.Fatalf("expected 1 drifted, got %d", len(r.Drifted))
	}
	if r.Drifted[0].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT to drift, got %s", r.Drifted[0].Key)
	}
	if r.Drifted[0].Pinned != "8080" {
		t.Errorf("expected pinned=8080, got %s", r.Drifted[0].Pinned)
	}
	if r.Drifted[0].Current != "9090" {
		t.Errorf("expected current=9090, got %s", r.Drifted[0].Current)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := baseEnv()
	pins := PinAll(env)
	delete(env, "DB_HOST")
	r := Check(env, pins)
	var found bool
	for _, d := range r.Drifted {
		if d.Key == "DB_HOST" && d.Current == "<missing>" {
			found = true
		}
	}
	if !found {
		t.Error("expected DB_HOST to appear as missing drift")
	}
}

func TestHasDrift(t *testing.T) {
	r := Result{Drifted: []Drift{{Key: "X", Pinned: "a", Current: "b"}}}
	if !HasDrift(r) {
		t.Error("expected HasDrift=true")
	}
	if HasDrift(Result{}) {
		t.Error("expected HasDrift=false for empty result")
	}
}

func TestWriteReport_NoPins(t *testing.T) {
	var buf bytes.Buffer
	WriteReport(&buf, Result{})
	if !strings.Contains(buf.String(), "No pins") {
		t.Errorf("expected 'No pins' message, got: %s", buf.String())
	}
}

func TestWriteReport_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := Result{
		Drifted: []Drift{{Key: "APP_ENV", Pinned: "production", Current: "staging"}},
		Matched: []string{"APP_PORT"},
	}
	WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in report")
	}
	if !strings.Contains(out, "1 matched, 1 drifted") {
		t.Errorf("unexpected summary: %s", out)
	}
}
