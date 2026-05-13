package templater_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlint/templater"
)

func TestRender_NoPlaceholders(t *testing.T) {
	result := templater.Render("APP_ENV=production\nDEBUG=false", map[string]string{})
	if result.Rendered != "APP_ENV=production\nDEBUG=false" {
		t.Errorf("unexpected rendered output: %q", result.Rendered)
	}
	if len(result.Substituted) != 0 {
		t.Errorf("expected no substitutions, got %v", result.Substituted)
	}
}

func TestRender_SubstitutesKnownKeys(t *testing.T) {
	tmpl := "DB_HOST={{ DB_HOST }}\nDB_PORT={{ DB_PORT }}"
	values := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	result := templater.Render(tmpl, values)

	if result.Rendered != "DB_HOST=localhost\nDB_PORT=5432" {
		t.Errorf("unexpected rendered output: %q", result.Rendered)
	}
	if len(result.Substituted) != 2 {
		t.Errorf("expected 2 substitutions, got %d", len(result.Substituted))
	}
	if len(result.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", result.Unresolved)
	}
}

func TestRender_TracksUnresolved(t *testing.T) {
	tmpl := "SECRET={{ SECRET }}\nHOST={{ HOST }}"
	result := templater.Render(tmpl, map[string]string{"HOST": "localhost"})

	if len(result.Unresolved) != 1 || result.Unresolved[0] != "SECRET" {
		t.Errorf("expected [SECRET] unresolved, got %v", result.Unresolved)
	}
	if len(result.Substituted) != 1 || result.Substituted[0] != "HOST" {
		t.Errorf("expected [HOST] substituted, got %v", result.Substituted)
	}
}

func TestRender_DeduplicatesUnresolved(t *testing.T) {
	tmpl := "A={{ MISSING }}\nB={{ MISSING }}"
	result := templater.Render(tmpl, map[string]string{})
	if len(result.Unresolved) != 1 {
		t.Errorf("expected 1 unique unresolved entry, got %d", len(result.Unresolved))
	}
}

func TestRenderFile_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "template.env")
	_ = os.WriteFile(path, []byte("PORT={{ PORT }}"), 0644)

	result, err := templater.RenderFile(path, map[string]string{"PORT": "8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Rendered != "PORT=8080" {
		t.Errorf("unexpected rendered output: %q", result.Rendered)
	}
}

func TestRenderFile_MissingFile(t *testing.T) {
	_, err := templater.RenderFile("/nonexistent/path.env", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
