package converter_test

import (
	"errors"
	"strings"
	"testing"

	"envlint/converter"
)

func TestWriteReport_Success(t *testing.T) {
	var sb strings.Builder
	converter.WriteReport(&sb, ".env", converter.FormatJSON, 5, nil)
	out := sb.String()
	if !strings.Contains(out, "✓") {
		t.Errorf("expected success icon in output, got: %s", out)
	}
	if !strings.Contains(out, "5 variables") {
		t.Errorf("expected variable count in output, got: %s", out)
	}
	if !strings.Contains(out, "json") {
		t.Errorf("expected format name in output, got: %s", out)
	}
}

func TestWriteReport_Failure(t *testing.T) {
	var sb strings.Builder
	converter.WriteReport(&sb, ".env", converter.Format("xml"), 0, errors.New("unsupported format"))
	out := sb.String()
	if !strings.Contains(out, "✗") {
		t.Errorf("expected failure icon in output, got: %s", out)
	}
	if !strings.Contains(out, "unsupported format") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := converter.SupportedFormats()
	if len(formats) == 0 {
		t.Error("expected at least one supported format")
	}
}

func TestFormatNames(t *testing.T) {
	names := converter.FormatNames()
	found := false
	for _, n := range names {
		if n == "json" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'json' in format names, got: %v", names)
	}
}
