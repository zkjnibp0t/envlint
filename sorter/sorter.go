// Package sorter provides utilities for sorting and grouping .env file entries
// alphabetically or by custom prefix groups.
package sorter

import (
	"sort"
	"strings"
)

// Result holds the sorted environment map and any detected groups.
type Result struct {
	Sorted map[string]string
	Groups map[string][]string
}

// Sort returns a new map with keys ordered alphabetically. The original map
// is never mutated.
func Sort(env map[string]string) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]string, len(env))
	for _, k := range keys {
		sorted[k] = env[k]
	}

	return Result{
		Sorted: sorted,
		Groups: groupByPrefix(keys),
	}
}

// groupByPrefix clusters keys by their first underscore-delimited segment.
// Keys without an underscore are placed under the "_" bucket.
func groupByPrefix(keys []string) map[string][]string {
	groups := make(map[string][]string)
	for _, k := range keys {
		prefix := prefix(k)
		groups[prefix] = append(groups[prefix], k)
	}
	return groups
}

func prefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return "_"
}

// SortedKeys returns the keys of env in alphabetical order.
func SortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
