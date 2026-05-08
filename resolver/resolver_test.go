package resolver_test

import (
	"testing"

	"github.com/user/envlint/resolver"
)

func sources(pairs ...interface{}) []resolver.NamedSource {
	var out []resolver.NamedSource
	for i := 0; i+1 < len(pairs); i += 2 {
		name := pairs[i].(string)
		env := pairs[i+1].(map[string]string)
		out = append(out, resolver.NamedSource{Name: name, Env: env})
	}
	return out
}

func TestResolve_AllFound(t *testing.T) {
	srcs := sources(".env", map[string]string{"PORT": "8080", "HOST": "localhost"})
	result := resolver.Resolve([]string{"PORT", "HOST"}, srcs)

	if len(result.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", result.Missing)
	}
	if result.Resolutions[0].Value != "8080" {
		t.Errorf("expected PORT=8080, got %s", result.Resolutions[0].Value)
	}
}

func TestResolve_MissingKey(t *testing.T) {
	srcs := sources(".env", map[string]string{"PORT": "8080"})
	result := resolver.Resolve([]string{"PORT", "DB_URL"}, srcs)

	if len(result.Missing) != 1 || result.Missing[0] != "DB_URL" {
		t.Errorf("expected DB_URL missing, got %v", result.Missing)
	}
}

func TestResolve_FirstSourceWins(t *testing.T) {
	srcs := sources(
		".env.local", map[string]string{"PORT": "9090"},
		".env", map[string]string{"PORT": "8080"},
	)
	result := resolver.Resolve([]string{"PORT"}, srcs)

	if result.Resolutions[0].Value != "9090" {
		t.Errorf("expected first source to win, got %s", result.Resolutions[0].Value)
	}
	if result.Resolutions[0].Source != ".env.local" {
		t.Errorf("expected source .env.local, got %s", result.Resolutions[0].Source)
	}
}

func TestResolve_DeduplicatesKeys(t *testing.T) {
	srcs := sources(".env", map[string]string{"PORT": "8080"})
	result := resolver.Resolve([]string{"PORT", "PORT"}, srcs)

	if len(result.Resolutions) != 1 {
		t.Errorf("expected 1 resolution, got %d", len(result.Resolutions))
	}
}

func TestResolve_EmptySources(t *testing.T) {
	result := resolver.Resolve([]string{"PORT"}, nil)

	if len(result.Missing) != 1 {
		t.Errorf("expected PORT to be missing, got %v", result.Missing)
	}
}

func TestResult_Summary_AllResolved(t *testing.T) {
	srcs := sources(".env", map[string]string{"A": "1", "B": "2"})
	result := resolver.Resolve([]string{"A", "B"}, srcs)

	got := result.Summary()
	if got != "all 2 variable(s) resolved" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestResult_Summary_WithMissing(t *testing.T) {
	srcs := sources(".env", map[string]string{"A": "1"})
	result := resolver.Resolve([]string{"A", "B"}, srcs)

	got := result.Summary()
	expected := "1/2 variable(s) resolved, 1 missing: [B]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
