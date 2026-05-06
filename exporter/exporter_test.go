package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/exporter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"APP_PORT": "8080",
	}
}

func TestExport_ShellFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv(), exporter.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected shell export for APP_ENV, got:\n%s", out)
	}
	if !strings.Contains(out, "export APP_PORT=8080") {
		t.Errorf("expected shell export for APP_PORT, got:\n%s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv(), exporter.FormatDocker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "export ") {
		t.Errorf("docker format should not contain 'export', got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in docker output, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv(), exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON to start with '{', got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV": "production"`) {
		t.Errorf("expected JSON entry for APP_ENV, got:\n%s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv(), exporter.FormatDocker)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// Sorted: APP_ENV, APP_PORT, DB_HOST
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected first line to be APP_ENV, got: %s", lines[0])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEnv(), exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_EmptyEnv(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, map[string]string{}, exporter.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty env, got: %s", buf.String())
	}
}
