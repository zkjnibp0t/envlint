package merger

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteReport_NoConflicts(t *testing.T) {
	r := Result{
		Env:       map[string]string{"FOO": "bar"},
		Conflicts: nil,
	}
	var buf bytes.Buffer
	WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "No conflicts") {
		t.Errorf("expected no-conflict message, got: %s", out)
	}
	if !strings.Contains(out, "1 keys") {
		t.Errorf("expected key count, got: %s", out)
	}
}

func TestWriteReport_WithConflicts(t *testing.T) {
	r := Result{
		Env: map[string]string{"DB_URL": "postgres://prod"},
		Conflicts: []Conflict{
			{
				Key:      "DB_URL",
				Winner:   "prod.env",
				Losers:   []string{"base.env"},
				FinalVal: "postgres://prod",
			},
		},
	}
	var buf bytes.Buffer
	WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "CONFLICT") {
		t.Errorf("expected CONFLICT label, got: %s", out)
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected key DB_URL, got: %s", out)
	}
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected winner name, got: %s", out)
	}
	if !strings.Contains(out, "base.env") {
		t.Errorf("expected loser name, got: %s", out)
	}
}

func TestWriteReport_SortedConflicts(t *testing.T) {
	r := Result{
		Env: map[string]string{},
		Conflicts: []Conflict{
			{Key: "Z_KEY", Winner: "a", FinalVal: "z"},
			{Key: "A_KEY", Winner: "b", FinalVal: "a"},
		},
	}
	var buf bytes.Buffer
	WriteReport(&buf, r)
	out := buf.String()
	azIdx := strings.Index(out, "A_KEY")
	zzIdx := strings.Index(out, "Z_KEY")
	if azIdx > zzIdx {
		t.Errorf("expected A_KEY before Z_KEY in output")
	}
}
