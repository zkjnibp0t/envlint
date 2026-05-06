package linter

import (
	"fmt"
	"io"

	"github.com/user/envlint/auditor"
	"github.com/user/envlint/formatter"
)

// WriteReport writes a human-readable summary of a lint Result to w.
func WriteReport(w io.Writer, r Result) {
	writeValidation(w, r)
	writeFormat(w, r)
	writeAudit(w, r)
	writeSuggestions(w, r)
	writeSummary(w, r)
}

func writeValidation(w io.Writer, r Result) {
	if len(r.ValidationErrors) == 0 {
		return
	}
	fmt.Fprintln(w, "=== Validation Errors ===")
	for _, e := range r.ValidationErrors {
		fmt.Fprintf(w, "  [ERROR] %s: %s\n", e.Key, e.Message)
	}
}

func writeFormat(w io.Writer, r Result) {
	if len(r.FormatIssues) == 0 {
		return
	}
	fmt.Fprintln(w, "=== Format Issues ===")
	for _, i := range r.FormatIssues {
		icon := formatter.LevelIcon(i.Level)
		fmt.Fprintf(w, "  %s line %d: %s\n", icon, i.Line, i.Message)
	}
}

func writeAudit(w io.Writer, r Result) {
	if len(r.AuditIssues) == 0 {
		return
	}
	fmt.Fprintln(w, "=== Audit Issues ===")
	for _, i := range r.AuditIssues {
		var level string
		switch i.Level {
		case auditor.LevelError:
			level = "ERROR"
		case auditor.LevelWarn:
			level = "WARN"
		default:
			level = "INFO"
		}
		fmt.Fprintf(w, "  [%s] %s: %s\n", level, i.Key, i.Message)
	}
}

func writeSuggestions(w io.Writer, r Result) {
	if len(r.Suggestions) == 0 {
		return
	}
	fmt.Fprintln(w, "=== Suggestions ===")
	for _, s := range r.Suggestions {
		fmt.Fprintf(w, "  Did you mean '%s' instead of '%s'?\n", s.Closest, s.Key)
	}
}

func writeSummary(w io.Writer, r Result) {
	fmt.Fprintln(w, "=== Summary ===")
	fmt.Fprintf(w, "  Validation errors : %d\n", len(r.ValidationErrors))
	fmt.Fprintf(w, "  Format issues     : %d\n", len(r.FormatIssues))
	fmt.Fprintf(w, "  Audit issues      : %d\n", len(r.AuditIssues))
	fmt.Fprintf(w, "  Suggestions       : %d\n", len(r.Suggestions))
	if r.HasErrors {
		fmt.Fprintln(w, "  Result            : FAIL")
	} else {
		fmt.Fprintln(w, "  Result            : PASS")
	}
}
