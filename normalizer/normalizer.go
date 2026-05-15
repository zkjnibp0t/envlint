// Package normalizer provides utilities to normalize .env file entries
// by standardizing key casing, trimming whitespace, and deduplicating entries.
package normalizer

import (
	"strings"
)

// Options controls normalization behaviour.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE when true.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// DeduplicateKeys keeps only the last occurrence of a duplicated key.
	DeduplicateKeys bool
}

// DefaultOptions returns a sensible default set of normalization options.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:   true,
		TrimValues:      true,
		DeduplicateKeys: true,
	}
}

// Result holds the normalized environment map and metadata.
type Result struct {
	Env        map[string]string
	Renamed    []string // keys that were uppercased
	Duplicates []string // keys that had duplicates removed
}

// Normalize applies the given options to the input env map and returns a Result.
func Normalize(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	seen := make(map[string]string) // normalized key -> original key

	var renamed []string
	var duplicates []string

	for k, v := range env {
		newKey := k
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(k)
			if newKey != k {
				renamed = append(renamed, k)
			}
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.DeduplicateKeys {
			if orig, exists := seen[newKey]; exists && orig != newKey {
				duplicates = append(duplicates, newKey)
			} else if _, exists := out[newKey]; exists {
				duplicates = append(duplicates, newKey)
			}
		}
		seen[newKey] = k
		out[newKey] = v
	}

	return Result{
		Env:        out,
		Renamed:    renamed,
		Duplicates: duplicates,
	}
}
