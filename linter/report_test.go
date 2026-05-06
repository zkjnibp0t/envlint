package linter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/linter"
	"github.com/user/envlint/validator"
)

func TestWriteReport_PassResult(t *testing.T) {
	r := linter.Result{HasErrors: false}
	var buf bytes.Buffer
	linter.WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS in output, got:\n%s", out)
	}
}

func TestWriteReport_FailResult(t *testing.T) {
	r := linter.Result{
		HasErrors: true,
		ValidationErrors: []validator.Error{
			{Key: "PORT", Message: "required variable is missing"},
		},
	}
	var buf bytes.Buffer
	linter.WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in validation section, got:\n%s", out)
	}
}

func TestWriteReport_SummaryCounts(t *testing.T) {
	r := linter.Result{
		ValidationErrors: []validator.Error{
			{Key: "X", Message: "bad"},
			{Key: "Y", Message: "also bad"},
		},
		HasErrors: true,
	}
	var buf bytes.Buffer
	linter.WriteReport(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Validation errors : 2") {
		t.Errorf("expected count 2 in summary, got:\n%s", out)
	}
}

func TestWriteReport_NoSectionsWhenClean(t *testing.T) {
	r := linter.Result{HasErrors: false}
	var buf bytes.Buffer
	linter.WriteReport(&buf, r)
	out := buf.String()
	for _, section := range []string{"Validation Errors", "Format Issues", "Audit Issues", "Suggestions"} {
		if strings.Contains(out, section) {
			t.Errorf("did not expect section %q in clean output", section)
		}
	}
}
