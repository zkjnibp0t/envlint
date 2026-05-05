package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlint/reporter"
	"github.com/user/envlint/validator"
)

func makeErrors(keys ...string) []validator.ValidationError {
	var errs []validator.ValidationError
	for _, k := range keys {
		errs = append(errs, validator.ValidationError{
			Key:      k,
			Message:  "is required but missing",
			Severity: validator.SeverityError,
		})
	}
	return errs
}

func TestWrite_TextNoErrors(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	ok := r.Write(nil)
	if !ok {
		t.Fatal("expected ok=true for no errors")
	}
	if !strings.Contains(buf.String(), "valid") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestWrite_TextWithErrors(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	ok := r.Write(makeErrors("DB_HOST", "API_KEY"))
	if ok {
		t.Fatal("expected ok=false when errors present")
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "2 validation error") {
		t.Errorf("expected error count in output, got: %s", out)
	}
}

func TestWrite_JSONNoErrors(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	r.Write(nil)
	out := buf.String()
	if !strings.Contains(out, `"valid":true`) {
		t.Errorf("expected valid:true, got: %s", out)
	}
	if !strings.Contains(out, `"errors":[]`) {
		t.Errorf("expected empty errors array, got: %s", out)
	}
}

func TestWrite_JSONWithErrors(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	r.Write(makeErrors("SECRET"))
	out := buf.String()
	if !strings.Contains(out, `"valid":false`) {
		t.Errorf("expected valid:false, got: %s", out)
	}
	if !strings.Contains(out, `"key":"SECRET"`) {
		t.Errorf("expected key SECRET, got: %s", out)
	}
}
