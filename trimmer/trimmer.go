// Package trimmer provides utilities for cleaning up .env files by removing
// redundant whitespace, duplicate keys, and blank lines from env maps and raw lines.
package trimmer

import (
	"strings"
)

// Result holds the output of a Trim operation.
type Result struct {
	// Env is the cleaned key-value map.
	Env map[string]string
	// RemovedDuplicates lists keys that appeared more than once (last value wins).
	RemovedDuplicates []string
	// BlankLinesRemoved is the count of blank/whitespace-only lines dropped.
	BlankLinesRemoved int
	// TrimmedValues lists keys whose values had surrounding whitespace stripped.
	TrimmedValues []string
}

// Trim processes raw .env file lines and returns a cleaned Result.
// Duplicate keys are collapsed (last value wins). Values are whitespace-trimmed.
func Trim(lines []string) Result {
	env := make(map[string]string)
	order := []string{}
	seen := make(map[string]bool)
	duplicates := make(map[string]bool)
	var blankRemoved int
	var trimmedValues []string

	for _, line := range lines {
		stripped := strings.TrimSpace(line)

		if stripped == "" || strings.HasPrefix(stripped, "#") {
			if stripped == "" {
				blankRemoved++
			}
			continue
		}

		idx := strings.IndexByte(stripped, '=')
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(stripped[:idx])
		val := strings.TrimSpace(stripped[idx+1:])

		if seen[key] {
			duplicates[key] = true
		} else {
			order = append(order, key)
			seen[key] = true
		}

		origVal := stripped[idx+1:]
		if origVal != val {
			trimmedValues = append(trimmedValues, key)
		}

		env[key] = val
	}

	var removedDups []string
	for k := range duplicates {
		removedDups = append(removedDups, k)
	}

	_ = order

	return Result{
		Env:               env,
		RemovedDuplicates: removedDups,
		BlankLinesRemoved: blankRemoved,
		TrimmedValues:     trimmedValues,
	}
}
