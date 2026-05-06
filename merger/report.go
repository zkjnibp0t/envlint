package merger

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a human-readable merge summary to w.
func WriteReport(w io.Writer, r Result) {
	fmt.Fprintf(w, "Merged %d keys\n", len(r.Env))

	if len(r.Conflicts) == 0 {
		fmt.Fprintln(w, "No conflicts detected.")
		return
	}

	fmt.Fprintf(w, "%d conflict(s) found:\n", len(r.Conflicts))

	// Sort for deterministic output.
	sorted := make([]Conflict, len(r.Conflicts))
	copy(sorted, r.Conflicts)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, c := range sorted {
		fmt.Fprintf(w, "  [CONFLICT] %s\n", c.Key)
		fmt.Fprintf(w, "    winner : %s (value=%q)\n", c.Winner, c.FinalVal)
		for _, l := range c.Losers {
			fmt.Fprintf(w, "    ignored: %s\n", l)
		}
	}
}
