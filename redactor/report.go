package redactor

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a human-readable summary of which keys were redacted.
func WriteReport(w io.Writer, original, redacted map[string]string) {
	var sensitive []string
	for k, v := range redacted {
		if v != original[k] {
			sensitive = append(sensitive, k)
		}
	}
	sort.Strings(sensitive)

	if len(sensitive) == 0 {
		fmt.Fprintln(w, "redactor: no sensitive keys found")
		return
	}

	fmt.Fprintf(w, "redactor: %d sensitive key(s) redacted\n", len(sensitive))
	for _, k := range sensitive {
		fmt.Fprintf(w, "  - %s\n", k)
	}
}
