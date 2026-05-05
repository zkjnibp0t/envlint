package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/formatter"
)

func makeIssues() []formatter.Issue {
	return []formatter.Issue{
		{Line: 1, Key: "db_host", Message: `key "db_host" should be uppercase`, Level: formatter.LevelWarning},
		{Line: 3, Key: "PORT", Message: `key "PORT" has trailing whitespace in value`, Level: formatter.LevelInfo},
	}
}

func TestWriteIssues_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	formatter.WriteIssues(&buf, nil)
	if !strings.Contains(buf.String(), "No style issues") {
		t.Errorf("expected no-issues message, got: %s", buf.String())
	}
}

func TestWriteIssues_WithIssues(t *testing.T) {
	var buf bytes.Buffer
	formatter.WriteIssues(&buf, makeIssues())
	out := buf.String()
	if !strings.Contains(out, "2 style issue") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected WARNING label, got: %s", out)
	}
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO label, got: %s", out)
	}
}

func TestSummary_NoIssues(t *testing.T) {
	s := formatter.Summary(nil)
	if s != "style: ok" {
		t.Errorf("expected 'style: ok', got %q", s)
	}
}

func TestSummary_WithIssues(t *testing.T) {
	s := formatter.Summary(makeIssues())
	if !strings.Contains(s, "warning") {
		t.Errorf("expected warning in summary, got %q", s)
	}
	if !strings.Contains(s, "info") {
		t.Errorf("expected info in summary, got %q", s)
	}
}
