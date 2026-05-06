package converter_test

import (
	"strings"
	"testing"

	"envlint/converter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_PORT":  "5432",
		"API_KEY":  "secret123",
	}
}

func TestConvert_JSON(t *testing.T) {
	out, err := converter.Convert(sampleEnv(), converter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected APP_ENV in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"production"`) {
		t.Errorf("expected production value in JSON output")
	}
}

func TestConvert_YAML(t *testing.T) {
	out, err := converter.Convert(sampleEnv(), converter.FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV:") {
		t.Errorf("expected APP_ENV key in YAML output, got:\n%s", out)
	}
}

func TestConvert_TOML(t *testing.T) {
	out, err := converter.Convert(sampleEnv(), converter.FormatTOML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_PORT = ") {
		t.Errorf("expected DB_PORT key in TOML output, got:\n%s", out)
	}
}

func TestConvert_Shell(t *testing.T) {
	out, err := converter.Convert(sampleEnv(), converter.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export API_KEY=") {
		t.Errorf("expected export statement in shell output, got:\n%s", out)
	}
}

func TestConvert_DotEnv(t *testing.T) {
	out, err := converter.Convert(sampleEnv(), converter.FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in dotenv output, got:\n%s", out)
	}
}

func TestConvert_SortedOutput(t *testing.T) {
	out, err := converter.Convert(sampleEnv(), converter.FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "API_KEY") {
		t.Errorf("expected sorted output to start with API_KEY, got: %s", lines[0])
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	_, err := converter.Convert(sampleEnv(), converter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
