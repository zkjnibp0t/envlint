package scanner

import (
	"fmt"
	"io"
	"sort"
)

const (
	iconHigh   = "🔴"
	iconMedium = "🟡"
	iconLow    = "🔵"
)

// WriteReport writes a human-readable security scan report to w.
func WriteReport(w io.Writer, findings []Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "✅ No security issues detected.")
		return
	}

	// Sort by severity then key for deterministic output.
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			return severityOrder(findings[i].Severity) < severityOrder(findings[j].Severity)
		}
		return findings[i].Key < findings[j].Key
	})

	fmt.Fprintf(w, "⚠️  Security scan found %d issue(s):\n\n", len(findings))

	for _, f := range findings {
		icon := iconForSeverity(f.Severity)
		fmt.Fprintf(w, "  %s [%s] %s\n", icon, f.Severity, f.Key)
		fmt.Fprintf(w, "     rule    : %s\n", f.Rule)
		fmt.Fprintf(w, "     message : %s\n", f.Message)
		if f.Value != "" && f.Value != "[REDACTED]" {
			fmt.Fprintf(w, "     value   : %q\n", f.Value)
		}
		fmt.Fprintln(w)
	}

	high, medium, low := countBySeverity(findings)
	fmt.Fprintf(w, "Summary: %d high, %d medium, %d low\n", high, medium, low)
}

func iconForSeverity(s Severity) string {
	switch s {
	case SeverityHigh:
		return iconHigh
	case SeverityMedium:
		return iconMedium
	default:
		return iconLow
	}
}

func severityOrder(s Severity) int {
	switch s {
	case SeverityHigh:
		return 0
	case SeverityMedium:
		return 1
	default:
		return 2
	}
}

func countBySeverity(findings []Finding) (high, medium, low int) {
	for _, f := range findings {
		switch f.Severity {
		case SeverityHigh:
			high++
		case SeverityMedium:
			medium++
		default:
			low++
		}
	}
	return
}
