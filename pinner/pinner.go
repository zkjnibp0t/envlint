// Package pinner provides utilities for pinning environment variable values
// to a snapshot, detecting drift between the current env and a previously
// recorded baseline.
package pinner

import (
	"fmt"
	"sort"
	"time"
)

// Pin represents a recorded value for a single environment variable.
type Pin struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// Result holds the outcome of comparing the current env against pinned values.
type Result struct {
	Drifted []Drift
	Matched []string
}

// Drift describes a variable whose current value differs from its pinned value.
type Drift struct {
	Key      string
	Pinned   string
	Current  string
}

// PinAll records the current value of every key in env as a Pin.
func PinAll(env map[string]string) []Pin {
	now := time.Now().UTC()
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pins := make([]Pin, 0, len(keys))
	for _, k := range keys {
		pins = append(pins, Pin{Key: k, Value: env[k], PinnedAt: now})
	}
	return pins
}

// Check compares the current env against a slice of Pins and returns a Result
// describing which keys have drifted and which still match.
func Check(env map[string]string, pins []Pin) Result {
	var result Result
	for _, p := range pins {
		current, ok := env[p.Key]
		if !ok {
			result.Drifted = append(result.Drifted, Drift{
				Key:     p.Key,
				Pinned:  p.Value,
				Current: "<missing>",
			})
			continue
		}
		if current != p.Value {
			result.Drifted = append(result.Drifted, Drift{
				Key:     p.Key,
				Pinned:  p.Value,
				Current: current,
			})
		} else {
			result.Matched = append(result.Matched, p.Key)
		}
	}
	return result
}

// HasDrift reports whether any variables have drifted from their pinned values.
func HasDrift(r Result) bool {
	return len(r.Drifted) > 0
}

// Summary returns a human-readable one-line summary of the result.
func Summary(r Result) string {
	return fmt.Sprintf("%d matched, %d drifted", len(r.Matched), len(r.Drifted))
}
