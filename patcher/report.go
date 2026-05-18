package patcher

import (
	"fmt"
	"io"
)

func iconForOp(op Op) string {
	switch op {
	case OpAdded:
		return "+"
	case OpUpdated:
		return "~"
	default:
		return " "
	}
}

// WriteReport writes a human-readable summary of patch changes to w.
// Unchanged keys are omitted unless there are no changes at all.
func WriteReport(w io.Writer, result Result) {
	added, updated, unchanged := 0, 0, 0
	for _, c := range result.Changes {
		switch c.Op {
		case OpAdded:
			added++
		case OpUpdated:
			updated++
		default:
			unchanged++
		}
	}

	if added == 0 && updated == 0 {
		fmt.Fprintln(w, "patch: no changes applied")
		return
	}

	fmt.Fprintf(w, "patch: %d added, %d updated, %d unchanged\n", added, updated, unchanged)
	for _, c := range result.Changes {
		if c.Op == OpUnchanged {
			continue
		}
		if c.Op == OpAdded {
			fmt.Fprintf(w, "  [%s] %s = %q\n", iconForOp(c.Op), c.Key, c.NewValue)
		} else {
			fmt.Fprintf(w, "  [%s] %s: %q -> %q\n", iconForOp(c.Op), c.Key, c.OldValue, c.NewValue)
		}
	}
}
