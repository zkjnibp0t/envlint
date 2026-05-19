package pinner

import (
	"fmt"
	"io"
	"text/tabwriter"
)

const (
	iconOK    = "✔"
	iconDrift = "✘"
)

// WriteReport writes a formatted drift report to w.
func WriteReport(w io.Writer, r Result) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if len(r.Drifted) == 0 && len(r.Matched) == 0 {
		fmt.Fprintln(w, "No pins to check.")
		return
	}

	if len(r.Drifted) > 0 {
		fmt.Fprintln(w, "Drifted variables:")
		for _, d := range r.Drifted {
			if d.Current == "<missing>" {
				fmt.Fprintf(tw, "  %s\t%s\tpinned=%q\tcurrent=<missing>\n",
					iconDrift, d.Key, d.Pinned)
			} else {
				fmt.Fprintf(tw, "  %s\t%s\tpinned=%q\tcurrent=%q\n",
					iconDrift, d.Key, d.Pinned, d.Current)
			}
		}
		tw.Flush()
	}

	if len(r.Matched) > 0 {
		fmt.Fprintln(w, "Matched variables:")
		for _, k := range r.Matched {
			fmt.Fprintf(tw, "  %s\t%s\n", iconOK, k)
		}
		tw.Flush()
	}

	fmt.Fprintf(w, "\nSummary: %s\n", Summary(r))
}
