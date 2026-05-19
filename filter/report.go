package filter

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a human-readable summary of the filter result to w.
func WriteReport(w io.Writer, r Result) {
	kept := make([]string, len(r.Kept))
	copy(kept, r.Kept)
	sort.Strings(kept)

	dropped := make([]string, len(r.Dropped))
	copy(dropped, r.Dropped)
	sort.Strings(dropped)

	fmt.Fprintf(w, "Filter result: %d kept, %d dropped\n", len(kept), len(dropped))

	if len(kept) > 0 {
		fmt.Fprintln(w, "\nKept:")
		for _, k := range kept {
			fmt.Fprintf(w, "  ✔ %s = %s\n", k, r.Env[k])
		}
	}

	if len(dropped) > 0 {
		fmt.Fprintln(w, "\nDropped:")
		for _, k := range dropped {
			fmt.Fprintf(w, "  ✖ %s\n", k)
		}
	}
}
