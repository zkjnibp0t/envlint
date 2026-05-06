package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/differ"
)

func TestWriteReport_Clean(t *testing.T) {
	var buf bytes.Buffer
	d := differ.Diff{}
	differ.WriteReport(&buf, d, "a", "b")
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %q", buf.String())
	}
}

func TestWriteReport_OnlyInA(t *testing.T) {
	var buf bytes.Buffer
	d := differ.Diff{OnlyInA: []string{"SECRET"}}
	differ.WriteReport(&buf, d, ".env.example", ".env")
	out := buf.String()
	if !strings.Contains(out, ".env.example") {
		t.Errorf("expected label .env.example in output: %q", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in output: %q", out)
	}
}

func TestWriteReport_OnlyInB(t *testing.T) {
	var buf bytes.Buffer
	d := differ.Diff{OnlyInB: []string{"NEW_VAR"}}
	differ.WriteReport(&buf, d, "a", "b")
	out := buf.String()
	if !strings.Contains(out, "NEW_VAR") {
		t.Errorf("expected NEW_VAR in output: %q", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' marker for OnlyInB: %q", out)
	}
}

func TestWriteReport_Changed(t *testing.T) {
	var buf bytes.Buffer
	d := differ.Diff{Changed: []string{"PORT"}}
	differ.WriteReport(&buf, d, "a", "b")
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output: %q", out)
	}
	if !strings.Contains(out, "~") {
		t.Errorf("expected '~' marker for Changed: %q", out)
	}
}
