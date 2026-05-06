package interpolator_test

import (
	"strings"
	"testing"

	"envlint/interpolator"
)

func TestExpand_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	result, errs := interpolator.Expand(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if result["HOST"] != "localhost" || result["PORT"] != "5432" {
		t.Errorf("unexpected values: %v", result)
	}
}

func TestExpand_ResolvesInternalRef(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://localhost",
		"API_URL":  "${BASE_URL}/api",
	}
	result, errs := interpolator.Expand(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if result["API_URL"] != "http://localhost/api" {
		t.Errorf("expected 'http://localhost/api', got %q", result["API_URL"])
	}
}

func TestExpand_FallsBackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	env := map[string]string{
		"MY_VAR": "${OS_VAR}-suffix",
	}
	result, errs := interpolator.Expand(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if result["MY_VAR"] != "from-os-suffix" {
		t.Errorf("expected 'from-os-suffix', got %q", result["MY_VAR"])
	}
}

func TestExpand_UnresolvedRef(t *testing.T) {
	env := map[string]string{
		"DB_URL": "${MISSING_VAR}/db",
	}
	_, errs := interpolator.Expand(env)
	if len(errs) == 0 {
		t.Fatal("expected an error for unresolved reference")
	}
	if !strings.Contains(errs[0].Error(), "MISSING_VAR") {
		t.Errorf("error should mention MISSING_VAR, got: %v", errs[0])
	}
}

func TestExpand_MultipleRefsInValue(t *testing.T) {
	env := map[string]string{
		"PROTO":    "https",
		"DOMAIN":   "example.com",
		"FULL_URL": "${PROTO}://${DOMAIN}",
	}
	result, errs := interpolator.Expand(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if result["FULL_URL"] != "https://example.com" {
		t.Errorf("expected 'https://example.com', got %q", result["FULL_URL"])
	}
}
