// Package renamer provides utilities for bulk-renaming environment variable
// keys across a parsed .env map, applying a set of rename rules.
package renamer

import "fmt"

// Rule describes a single rename operation: from OldKey to NewKey.
type Rule struct {
	OldKey string
	NewKey string
}

// Result holds the outcome of applying all rename rules to an env map.
type Result struct {
	// Renamed contains rules that were successfully applied.
	Renamed []Rule
	// Skipped contains rules whose OldKey was not found in the env map.
	Skipped []Rule
	// Conflicts contains rules whose NewKey already existed in the env map
	// before the rename was applied (old key was left untouched).
	Conflicts []Rule
}

// Rename applies the given rules to env, returning a new map with keys
// renamed according to the rules. The original map is never mutated.
//
// Rules are applied in order. If a rule's OldKey does not exist the rule is
// recorded as Skipped. If a rule's NewKey already exists the rule is recorded
// as a Conflict and the old key is preserved.
func Rename(env map[string]string, rules []Rule) (map[string]string, Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var res Result
	for _, r := range rules {
		if r.OldKey == "" || r.NewKey == "" {
			res.Skipped = append(res.Skipped, r)
			continue
		}
		val, exists := out[r.OldKey]
		if !exists {
			res.Skipped = append(res.Skipped, r)
			continue
		}
		if _, conflict := out[r.NewKey]; conflict {
			res.Conflicts = append(res.Conflicts, r)
			continue
		}
		delete(out, r.OldKey)
		out[r.NewKey] = val
		res.Renamed = append(res.Renamed, r)
	}
	return out, res
}

// ParseRules converts a slice of "OLD=NEW" strings into Rule values.
// It returns an error if any entry is malformed.
func ParseRules(pairs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(pairs))
	for _, p := range pairs {
		for i, ch := range p {
			if ch == '=' {
				rules = append(rules, Rule{OldKey: p[:i], NewKey: p[i+1:]})
				goto next
			}
		}
		return nil, fmt.Errorf("renamer: invalid rule %q: expected OLD=NEW format", p)
	next:
	}
	return rules, nil
}
