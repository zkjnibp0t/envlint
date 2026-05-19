// Package filter provides utilities for selecting a subset of environment
// variables based on key patterns, prefixes, or custom predicates.
package filter

import (
	"regexp"
	"strings"
)

// Options controls how the filter is applied.
type Options struct {
	// Prefixes keeps only keys that start with any of the given prefixes.
	Prefixes []string

	// Pattern keeps only keys matching the regular expression.
	Pattern string

	// Exclude removes keys matching the regular expression after other filters.
	Exclude string

	// CaseSensitive controls whether prefix matching is case-sensitive.
	CaseSensitive bool
}

// Result holds the filtered environment and metadata.
type Result struct {
	Env     map[string]string
	Kept    []string
	Dropped []string
}

// Filter returns a new map containing only the key-value pairs from env that
// satisfy the given options. At least one option must be non-empty; if none
// are set every key is kept.
func Filter(env map[string]string, opts Options) (Result, error) {
	var includeRe, excludeRe *regexp.Regexp
	var err error

	if opts.Pattern != "" {
		includeRe, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, err
		}
	}

	if opts.Exclude != "" {
		excludeRe, err = regexp.Compile(opts.Exclude)
		if err != nil {
			return Result{}, err
		}
	}

	out := make(map[string]string)
	var kept, dropped []string

	for k, v := range env {
		if keep(k, opts, includeRe, excludeRe) {
			out[k] = v
			kept = append(kept, k)
		} else {
			dropped = append(dropped, k)
		}
	}

	return Result{Env: out, Kept: kept, Dropped: dropped}, nil
}

func keep(key string, opts Options, includeRe, excludeRe *regexp.Regexp) bool {
	cmp := key
	if !opts.CaseSensitive {
		cmp = strings.ToUpper(key)
	}

	if len(opts.Prefixes) > 0 {
		matched := false
		for _, p := range opts.Prefixes {
			pfx := p
			if !opts.CaseSensitive {
				pfx = strings.ToUpper(p)
			}
			if strings.HasPrefix(cmp, pfx) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if includeRe != nil && !includeRe.MatchString(key) {
		return false
	}

	if excludeRe != nil && excludeRe.MatchString(key) {
		return false
	}

	return true
}
