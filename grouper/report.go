package grouper

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a human-readable summary of a grouping Result to w.
func WriteReport(w io.Writer, r Result) {
	if len(r.Groups) == 0 && len(r.Ungrouped) == 0 {
		fmt.Fprintln(w, "no variables to group")
		return
	}

	for _, g := range r.Groups {
		fmt.Fprintf(w, "[%s] (%d keys)\n", g.Prefix, len(g.Keys))
		for _, k := range sortedKeys(g.Keys) {
			fmt.Fprintf(w, "  %s=%s\n", k, g.Keys[k])
		}
	}

	if len(r.Ungrouped) > 0 {
		fmt.Fprintf(w, "[UNGROUPED] (%d keys)\n", len(r.Ungrouped))
		for _, k := range sortedKeys(r.Ungrouped) {
			fmt.Fprintf(w, "  %s=%s\n", k, r.Ungrouped[k])
		}
	}

	total := len(r.Ungrouped)
	for _, g := range r.Groups {
		total += len(g.Keys)
	}
	fmt.Fprintf(w, "\ntotal: %d variable(s) across %d group(s)\n", total, len(r.Groups))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
