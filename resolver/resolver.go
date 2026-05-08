// Package resolver provides functionality to resolve and validate
// environment variable references across multiple .env files or sources.
package resolver

import (
	"fmt"
	"sort"
)

// Resolution holds the result of resolving a single variable.
type Resolution struct {
	Key    string
	Value  string
	Source string // which source file/map provided the value
	Found  bool
}

// Result holds all resolutions for a Resolve call.
type Result struct {
	Resolutions []Resolution
	Missing     []string
}

// Resolve looks up each key in the provided ordered list of sources.
// Sources are checked in order; the first one that contains a key wins.
// keys is the list of variable names to resolve.
// sources is an ordered slice of named maps: [{"name": map}, ...].
func Resolve(keys []string, sources []NamedSource) Result {
	seen := make(map[string]bool)
	resolutions := make([]Resolution, 0, len(keys))
	missing := []string{}

	for _, key := range keys {
		if seen[key] {
			continue
		}
		seen[key] = true

		resolved := false
		for _, src := range sources {
			if val, ok := src.Env[key]; ok {
				resolutions = append(resolutions, Resolution{
					Key:    key,
					Value:  val,
					Source: src.Name,
					Found:  true,
				})
				resolved = true
				break
			}
		}
		if !resolved {
			missing = append(missing, key)
			resolutions = append(resolutions, Resolution{
				Key:   key,
				Found: false,
			})
		}
	}

	sort.Strings(missing)
	return Result{Resolutions: resolutions, Missing: missing}
}

// NamedSource pairs a display name with an env map.
type NamedSource struct {
	Name string
	Env  map[string]string
}

// Summary returns a human-readable summary line for the result.
func (r Result) Summary() string {
	total := len(r.Resolutions)
	found := total - len(r.Missing)
	if len(r.Missing) == 0 {
		return fmt.Sprintf("all %d variable(s) resolved", total)
	}
	return fmt.Sprintf("%d/%d variable(s) resolved, %d missing: %v", found, total, len(r.Missing), r.Missing)
}
