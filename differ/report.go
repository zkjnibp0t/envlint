package differ

import (
	"fmt"
	"io"
)

// WriteReport writes a human-readable diff report to w.
// labelA and labelB are display names for the two env sources
// (e.g. ".env.example" and ".env").
func WriteReport(w io.Writer, d Diff, labelA, labelB string) {
	if d.IsClean() {
		fmt.Fprintln(w, "✔  No differences found.")
		return
	}

	if len(d.OnlyInA) > 0 {
		fmt.Fprintf(w, "Keys only in %s:\n", labelA)
		for _, k := range d.OnlyInA {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(d.OnlyInB) > 0 {
		fmt.Fprintf(w, "Keys only in %s:\n", labelB)
		for _, k := range d.OnlyInB {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(d.Changed) > 0 {
		fmt.Fprintln(w, "Keys with different values:")
		for _, k := range d.Changed {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}
}
