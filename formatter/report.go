package formatter

import (
	"fmt"
	"io"
	"strings"
)

// WriteIssues writes style issues to the provided writer in a human-readable format.
func WriteIssues(w io.Writer, issues []Issue) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No style issues found.")
		return
	}

	fmt.Fprintf(w, "Found %d style issue(s):\n", len(issues))
	for _, iss := range issues {
		icon := levelIcon(iss.Level)
		fmt.Fprintf(w, "  %s [line %d] %s: %s\n", icon, iss.Line, strings.ToUpper(string(iss.Level)), iss.Message)
	}
}

// Summary returns a compact string summary of issues grouped by level.
func Summary(issues []Issue) string {
	if len(issues) == 0 {
		return "style: ok"
	}
	counts := map[Level]int{}
	for _, iss := range issues {
		counts[iss.Level]++
	}
	parts := []string{}
	for _, lvl := range []Level{LevelError, LevelWarning, LevelInfo} {
		if n, ok := counts[lvl]; ok {
			parts = append(parts, fmt.Sprintf("%d %s(s)", n, lvl))
		}
	}
	return "style: " + strings.Join(parts, ", ")
}

func levelIcon(l Level) string {
	switch l {
	case LevelError:
		return "✖"
	case LevelWarning:
		return "⚠"
	case LevelInfo:
		return "ℹ"
	default:
		return "•"
	}
}
