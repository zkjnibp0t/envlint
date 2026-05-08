package scanner

import (
	"strings"
	"testing"
)

func TestScan_WeakPassword(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "password",
	}
	findings := Scan(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "weak-value" {
		t.Errorf("expected rule weak-value, got %s", findings[0].Rule)
	}
	if findings[0].Severity != SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", findings[0].Severity)
	}
}

func TestScan_EmbeddedPrivateKey(t *testing.T) {
	env := map[string]string{
		"SOME_KEY": "-----BEGIN RSA PRIVATE KEY-----\nMIIE...",
	}
	findings := Scan(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "embedded-private-key" {
		t.Errorf("expected rule embedded-private-key, got %s", findings[0].Rule)
	}
	if findings[0].Value != "[REDACTED]" {
		t.Errorf("expected value to be redacted")
	}
}

func TestScan_LocalhostURL(t *testing.T) {
	env := map[string]string{
		"DATABASE_HOST": "localhost",
	}
	findings := Scan(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "localhost-url" {
		t.Errorf("expected rule localhost-url, got %s", findings[0].Rule)
	}
	if findings[0].Severity != SeverityMedium {
		t.Errorf("expected MEDIUM severity")
	}
}

func TestScan_EmptySensitiveKey(t *testing.T) {
	env := map[string]string{
		"API_SECRET": "",
	}
	findings := Scan(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "empty-sensitive" {
		t.Errorf("expected rule empty-sensitive, got %s", findings[0].Rule)
	}
}

func TestScan_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	findings := Scan(env)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestWriteReport_NoFindings(t *testing.T) {
	var sb strings.Builder
	WriteReport(&sb, nil)
	if !strings.Contains(sb.String(), "No security issues") {
		t.Errorf("expected clean message, got: %s", sb.String())
	}
}

func TestWriteReport_WithFindings(t *testing.T) {
	findings := []Finding{
		{Key: "DB_PASSWORD", Value: "secret", Rule: "weak-value", Message: "weak", Severity: SeverityHigh},
		{Key: "DB_HOST", Value: "localhost", Rule: "localhost-url", Message: "localhost", Severity: SeverityMedium},
	}
	var sb strings.Builder
	WriteReport(&sb, findings)
	out := sb.String()
	if !strings.Contains(out, "HIGH") {
		t.Errorf("expected HIGH in output")
	}
	if !strings.Contains(out, "MEDIUM") {
		t.Errorf("expected MEDIUM in output")
	}
	if !strings.Contains(out, "Summary:") {
		t.Errorf("expected Summary line")
	}
}
