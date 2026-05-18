// Package patcher applies a set of key-value updates to an existing env map,
// tracking which keys were added, updated, or left unchanged.
package patcher

import "sort"

// Op describes the kind of change applied to a key.
type Op string

const (
	OpAdded    Op = "added"
	OpUpdated  Op = "updated"
	OpUnchanged Op = "unchanged"
)

// Change records what happened to a single key during a patch.
type Change struct {
	Key      string
	Op       Op
	OldValue string
	NewValue string
}

// Result holds the patched environment and the list of changes.
type Result struct {
	Env     map[string]string
	Changes []Change
}

// Patch applies the given patches to env and returns a new map along with
// a sorted list of Change records. The original env map is not mutated.
func Patch(env map[string]string, patches map[string]string) Result {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	changeMap := make(map[string]Change)

	// Mark all existing keys as unchanged initially.
	for k, v := range env {
		changeMap[k] = Change{Key: k, Op: OpUnchanged, OldValue: v, NewValue: v}
	}

	// Apply patches.
	for k, v := range patches {
		if old, exists := env[k]; exists {
			if old == v {
				changeMap[k] = Change{Key: k, Op: OpUnchanged, OldValue: old, NewValue: v}
			} else {
				changeMap[k] = Change{Key: k, Op: OpUpdated, OldValue: old, NewValue: v}
				out[k] = v
			}
		} else {
			changeMap[k] = Change{Key: k, Op: OpAdded, OldValue: "", NewValue: v}
			out[k] = v
		}
	}

	changes := make([]Change, 0, len(changeMap))
	for _, c := range changeMap {
		changes = append(changes, c)
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Result{Env: out, Changes: changes}
}
