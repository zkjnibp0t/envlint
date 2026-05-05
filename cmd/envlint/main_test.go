package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envlint")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return bin
}

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestMain_ValidEnv(t *testing.T) {
	bin := buildBinary(t)

	schemaPath := writeTempFile(t, "schema.yaml", `vars:
  PORT:
    type: int
    required: true
`)
	envPath := writeTempFile(t, ".env", "PORT=8080\n")

	cmd := exec.Command(bin, "--env", envPath, "--schema", schemaPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_MissingRequired(t *testing.T) {
	bin := buildBinary(t)

	schemaPath := writeTempFile(t, "schema.yaml", `vars:
  DATABASE_URL:
    type: url
    required: true
`)
	envPath := writeTempFile(t, ".env", "PORT=8080\n")

	cmd := exec.Command(bin, "--env", envPath, "--schema", schemaPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit code for missing required var, output: %s", out)
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	}
}

func TestMain_JSONFormat(t *testing.T) {
	bin := buildBinary(t)

	schemaPath := writeTempFile(t, "schema.yaml", `vars:
  PORT:
    type: int
    required: true
`)
	envPath := writeTempFile(t, ".env", "PORT=3000\n")

	cmd := exec.Command(bin, "--env", envPath, "--schema", schemaPath, "--format", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
}
