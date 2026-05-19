// Package grouper organises a flat env map into named groups based on
// key prefix, making it easier to reason about large .env files.
package grouper

import (
	"sort"
	"strings"
)

// Group holds all key/value pairs that share a common prefix.
type Group struct {
	Prefix string
	Keys   map[string]string
}

// Result is the output of a Group operation.
type Result struct {
	Groups    []Group
	Ungrouped map[string]string // keys with no recognised prefix
}

// ByPrefix splits env into groups using the provided prefixes.
// Each prefix is matched case-insensitively against the start of every key.
// Keys that match no prefix land in Result.Ungrouped.
// If prefixes is empty, all keys are placed in Ungrouped.
func ByPrefix(env map[string]string, prefixes []string) Result {
	groupMap := make(map[string]map[string]string, len(prefixes))
	for _, p := range prefixes {
		groupMap[strings.ToUpper(p)] = make(map[string]string)
	}

	ungrouped := make(map[string]string)

	for k, v := range env {
		matched := false
		for _, p := range prefixes {
			up := strings.ToUpper(p)
			if strings.HasPrefix(strings.ToUpper(k), up) {
				groupMap[up][k] = v
				matched = true
				break
			}
		}
		if !matched {
			ungrouped[k] = v
		}
	}

	groups := make([]Group, 0, len(prefixes))
	for _, p := range prefixes {
		up := strings.ToUpper(p)
		if len(groupMap[up]) > 0 {
			groups = append(groups, Group{Prefix: up, Keys: groupMap[up]})
		}
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})

	return Result{Groups: groups, Ungrouped: ungrouped}
}
