package renamer

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a human-readable summary of a rename Result to w.
func WriteReport(w io.Writer, res Result) {
	sort.Slice(res.Renamed, func(i, j int) bool { return res.Renamed[i].OldKey < res.Renamed[j].OldKey })
	sort.Slice(res.Skipped, func(i, j int) bool { return res.Skipped[i].OldKey < res.Skipped[j].OldKey })
	sort.Slice(res.Conflicts, func(i, j int) bool { return res.Conflicts[i].OldKey < res.Conflicts[j].OldKey })

	if len(res.Renamed) == 0 && len(res.Skipped) == 0 && len(res.Conflicts) == 0 {
		fmt.Fprintln(w, "✔ No rename rules provided.")
		return
	}

	if len(res.Renamed) > 0 {
		fmt.Fprintf(w, "✔ Renamed (%d):\n", len(res.Renamed))
		for _, r := range res.Renamed {
			fmt.Fprintf(w, "  %s → %s\n", r.OldKey, r.NewKey)
		}
	}

	if len(res.Conflicts) > 0 {
		fmt.Fprintf(w, "⚠ Conflicts (%d):\n", len(res.Conflicts))
		for _, r := range res.Conflicts {
			fmt.Fprintf(w, "  %s → %s  (target already exists, skipped)\n", r.OldKey, r.NewKey)
		}
	}

	if len(res.Skipped) > 0 {
		fmt.Fprintf(w, "✗ Skipped (%d):\n", len(res.Skipped))
		for _, r := range res.Skipped {
			fmt.Fprintf(w, "  %s (not found)\n", r.OldKey)
		}
	}

	fmt.Fprintf(w, "\nSummary: %d renamed, %d conflicts, %d skipped.\n",
		len(res.Renamed), len(res.Conflicts), len(res.Skipped))
}
