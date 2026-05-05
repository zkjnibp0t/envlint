package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envlint/validator"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes validation results to an output stream.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a Reporter writing to w with the given format.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w, format: format}
}

// Write outputs the validation errors. Returns true if there were no errors.
func (r *Reporter) Write(errs []validator.ValidationError) bool {
	if r.format == FormatJSON {
		r.writeJSON(errs)
	} else {
		r.writeText(errs)
	}
	return len(errs) == 0
}

func (r *Reporter) writeText(errs []validator.ValidationError) {
	if len(errs) == 0 {
		fmt.Fprintln(r.w, "✓ All environment variables are valid.")
		return
	}
	fmt.Fprintf(r.w, "✗ Found %d validation error(s):\n", len(errs))
	for _, e := range errs {
		fmt.Fprintf(r.w, "  - [%s] %s: %s\n", strings.ToUpper(string(e.Severity)), e.Key, e.Message)
	}
}

func (r *Reporter) writeJSON(errs []validator.ValidationError) {
	if len(errs) == 0 {
		fmt.Fprintln(r.w, `{"valid":true,"errors":[]}`)
		return
	}
	fmt.Fprintf(r.w, `{"valid":false,"errors":[`)
	for i, e := range errs {
		if i > 0 {
			fmt.Fprint(r.w, ",")
		}
		fmt.Fprintf(r.w, `{"key":%q,"message":%q,"severity":%q}`, e.Key, e.Message, e.Severity)
	}
	fmt.Fprintln(r.w, `]}`)
}
